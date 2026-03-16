package main

import (
    "bufio"
    "encoding/hex"
    "errors"
    "flag"
    "fmt"
    "io"
    "net"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "sync/atomic"
    "syscall"
    "time"
)

type Config struct {
    mode    string
    addr    string
    outHex  bool
    outStr  bool
    bufSize int
}

type ConnStats struct {
    id        uint64
    startedAt time.Time
    lastAt    time.Time
    inBytes   uint64
    outBytes  uint64
}

var nextConnID uint64

func main() {
    var cfg Config
    flag.StringVar(&cfg.mode, "mode", "server", "server or client")
    flag.StringVar(&cfg.addr, "addr", ":9000", "listen/connect address")
    flag.BoolVar(&cfg.outHex, "hex", false, "output hex for received/sent data")
    flag.BoolVar(&cfg.outStr, "str", true, "output string for received/sent data")
    flag.IntVar(&cfg.bufSize, "buf", 4096, "read buffer size")
    flag.Parse()

    cfg.mode = strings.ToLower(strings.TrimSpace(cfg.mode))
    if cfg.bufSize <= 0 {
        cfg.bufSize = 4096
    }
    if !cfg.outHex && !cfg.outStr {
        cfg.outStr = true
    }

    switch cfg.mode {
    case "server":
        runServer(cfg)
    case "client":
        runClient(cfg)
    default:
        logf(0, "ERROR", "unknown mode: %s", cfg.mode)
        os.Exit(2)
    }
}

func runServer(cfg Config) {
    ln, err := net.Listen("tcp", cfg.addr)
    if err != nil {
        logf(0, "ERROR", "listen failed: %v", err)
        os.Exit(1)
    }
    logf(0, "STATE", "LISTENING addr=%s", cfg.addr)

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigCh
        logf(0, "STATE", "SHUTDOWN signal received")
        _ = ln.Close()
    }()

    for {
        conn, err := ln.Accept()
        if err != nil {
            if ne, ok := err.(net.Error); ok && ne.Temporary() {
                logf(0, "WARN", "accept temporary error: %v", err)
                time.Sleep(200 * time.Millisecond)
                continue
            }
            logf(0, "STATE", "listener closed: %v", err)
            return
        }
        id := atomic.AddUint64(&nextConnID, 1)
        stats := &ConnStats{id: id, startedAt: time.Now(), lastAt: time.Now()}
        logf(id, "STATE", "ACCEPTED local=%s remote=%s", conn.LocalAddr(), conn.RemoteAddr())
        go handleConn(conn, cfg, stats, false)
    }
}

func runClient(cfg Config) {
    logf(0, "STATE", "CONNECTING addr=%s", cfg.addr)
    conn, err := net.Dial("tcp", cfg.addr)
    if err != nil {
        logf(0, "ERROR", "connect failed: %v", err)
        os.Exit(1)
    }
    id := atomic.AddUint64(&nextConnID, 1)
    stats := &ConnStats{id: id, startedAt: time.Now(), lastAt: time.Now()}
    logf(id, "STATE", "CONNECTED local=%s remote=%s", conn.LocalAddr(), conn.RemoteAddr())

    done := make(chan struct{})
    go func() {
        handleConn(conn, cfg, stats, true)
        close(done)
    }()

    // Read stdin and send to server
    stdin := bufio.NewScanner(os.Stdin)
    for stdin.Scan() {
        line := stdin.Text()
        if len(line) == 0 {
            continue
        }
        data := []byte(line + "\n")
        n, err := conn.Write(data)
        atomic.AddUint64(&stats.outBytes, uint64(n))
        stats.lastAt = time.Now()
        logData(stats.id, "OUT", data[:n], cfg)
        if err != nil {
            logf(stats.id, "ERROR", "write failed: %s", classifyClose(err))
            _ = conn.Close()
            <-done
            return
        }
    }
    if err := stdin.Err(); err != nil {
        logf(stats.id, "WARN", "stdin error: %v", err)
    }

    // half-close if possible
    if tcp, ok := conn.(*net.TCPConn); ok {
        _ = tcp.CloseWrite()
    } else {
        _ = conn.Close()
    }
    <-done
}

func handleConn(conn net.Conn, cfg Config, stats *ConnStats, isClient bool) {
    defer func() {
        _ = conn.Close()
        duration := time.Since(stats.startedAt)
        logf(stats.id, "LIFE", "CLOSED duration=%s in=%d out=%d", duration.Truncate(time.Millisecond), atomic.LoadUint64(&stats.inBytes), atomic.LoadUint64(&stats.outBytes))
    }()

    logf(stats.id, "STATE", "ESTABLISHED")

    buf := make([]byte, cfg.bufSize)
    for {
        n, err := conn.Read(buf)
        if n > 0 {
            atomic.AddUint64(&stats.inBytes, uint64(n))
            stats.lastAt = time.Now()
            logData(stats.id, "IN", buf[:n], cfg)
        }
        if err != nil {
            reason := classifyClose(err)
            logf(stats.id, "STATE", "CLOSING reason=%s", reason)
            return
        }
        if isClient {
            // client read loop only, writer handles outgoing
            continue
        }
    }
}

func logData(id uint64, dir string, data []byte, cfg Config) {
    if cfg.outStr {
        logf(id, "DATA", "%s STR %s", dir, safeString(data))
    }
    if cfg.outHex {
        logf(id, "DATA", "%s HEX %s", dir, hex.EncodeToString(data))
    }
}

func safeString(b []byte) string {
    s := string(b)
    // Quote non-printable bytes for reliable logs
    if isMostlyPrintable(s) {
        return s
    }
    return strconv.QuoteToASCII(s)
}

func isMostlyPrintable(s string) bool {
    printable := 0
    for _, r := range s {
        if r == '\n' || r == '\r' || r == '\t' || strconv.IsPrint(r) {
            printable++
        }
    }
    return printable*100 >= len([]rune(s))*80
}

func classifyClose(err error) string {
    if err == nil {
        return "unknown"
    }
    if errors.Is(err, io.EOF) {
        return "remote closed (FIN)"
    }
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Timeout() {
        return "timeout"
    }
    var opErr *net.OpError
    if errors.As(err, &opErr) {
        if opErr.Err != nil {
            if reason := syscallReason(opErr.Err); reason != "" {
                return reason
            }
        }
        return opErr.Err.Error()
    }
    if reason := syscallReason(err); reason != "" {
        return reason
    }
    return err.Error()
}

func syscallReason(err error) string {
    var se syscall.Errno
    if errors.As(err, &se) {
        switch se {
        case syscall.ECONNRESET:
            return "connection reset by peer (RST)"
        case syscall.EPIPE:
            return "broken pipe (write to closed connection)"
        case syscall.ETIMEDOUT:
            return "timeout"
        case syscall.ECONNABORTED:
            return "connection aborted"
        case syscall.ECONNREFUSED:
            return "connection refused"
        }
    }
    // Windows specific errno values may be wrapped differently
    msg := strings.ToLower(err.Error())
    if strings.Contains(msg, "wsarecv") && strings.Contains(msg, "10054") {
        return "connection reset by peer (RST)"
    }
    if strings.Contains(msg, "10053") {
        return "connection aborted"
    }
    if strings.Contains(msg, "10061") {
        return "connection refused"
    }
    return ""
}

func logf(id uint64, tag string, format string, args ...any) {
    ts := time.Now().Format("2006-01-02 15:04:05.000")
    prefix := fmt.Sprintf("%s [%s]", ts, tag)
    if id > 0 {
        prefix = fmt.Sprintf("%s [conn=%d]", prefix, id)
    }
    msg := fmt.Sprintf(format, args...)
    fmt.Printf("%s %s\n", prefix, msg)
}
