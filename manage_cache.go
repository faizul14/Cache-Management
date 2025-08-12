package main

import (
    "bufio"
    "fmt"
    "io/fs"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
    "time"
)

const (
    minSize    = 1024 // 1 KB minimal ukuran cache yang ditampilkan
    cacheDir   = ".cache"
    logFile    = "cache_management.log"
)

type CacheEntry struct {
    Path string
    Size int64
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    for {
        clearScreen()
        fmt.Println("=====================================")
        fmt.Println("        Cache Management by fmp")
        fmt.Println("=====================================")
        fmt.Println("1. List Cache")
        fmt.Println("2. Manage Cache")
        fmt.Println("3. Exit Program")
        fmt.Println("4. Tampilkan Info Disk Usage")
        fmt.Println("5. Log Aktivitas")
        fmt.Println("6. Bantuan")
        fmt.Println("=====================================")
        fmt.Print("Pilih menu (1-6): ")

        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        switch input {
        case "1":
            listCache()
        case "2":
            manageCache(reader)
        case "3":
            fmt.Println("👋 Keluar dari program...")
            os.Exit(0)
        case "4":
            diskInfo()
        case "5":
            showLog()
        case "6":
            helpMenu()
        default:
            fmt.Println("⚠ Pilihan tidak valid.")
        }

        fmt.Println()
        fmt.Print("Tekan ENTER untuk kembali ke menu...")
        reader.ReadString('\n')
    }
}

func clearScreen() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

func listCache() {
    cachePath := filepath.Join(os.Getenv("HOME"), cacheDir)
    entries, err := getCacheEntries(cachePath)
    if err != nil {
        fmt.Println("❌ Gagal membaca cache:", err)
        return
    }
    fmt.Println("=====================================")
    fmt.Printf("         LIST CACHE > %d bytes\n", minSize)
    fmt.Println("=====================================")
    fmt.Printf("%-4s | %-10s | %s\n", "No", "Size", "Folder")
    fmt.Println("----+------------+-------------------------")
    for i, e := range entries {
        fmt.Printf("%-4d | %-10s | %s\n", i+1, byteCountDecimal(e.Size), e.Path)
    }
    fmt.Println("----+------------+-------------------------")

    total := totalSize(entries)
    fmt.Printf("💾 Total Cache: %s\n", byteCountDecimal(total))
    fmt.Println("=====================================")
}

func manageCache(reader *bufio.Reader) {
    cachePath := filepath.Join(os.Getenv("HOME"), cacheDir)
    entries, err := getCacheEntries(cachePath)
    if err != nil {
        fmt.Println("❌ Gagal membaca cache:", err)
        return
    }
    if len(entries) == 0 {
        fmt.Println("✅ Tidak ada cache di atas batas minimum.")
        return
    }

    fmt.Println("=====================================")
    fmt.Printf("      MANAGE CACHE > %d bytes\n", minSize)
    fmt.Println("=====================================")
    fmt.Printf("%-4s | %-10s | %s\n", "No", "Size", "Folder")
    fmt.Println("----+------------+-------------------------")
    for i, e := range entries {
        fmt.Printf("%-4d | %-10s | %s\n", i+1, byteCountDecimal(e.Size), e.Path)
    }
    fmt.Println("----+------------+-------------------------")

    fmt.Println("📌 Ketik nomor cache yang ingin dihapus (pisahkan spasi)")
    fmt.Println("📌 Ketik Q untuk kembali ke menu utama")
    fmt.Print("Masukkan pilihan: ")
    input, _ := reader.ReadString('\n')
    input = strings.TrimSpace(input)
    if strings.EqualFold(input, "Q") {
        fmt.Println("🔙 Kembali ke menu utama...")
        return
    }

    selections := strings.Fields(input)
    for _, sel := range selections {
        idx, err := strconv.Atoi(sel)
        if err != nil || idx < 1 || idx > len(entries) {
            fmt.Printf("⚠ Nomor tidak valid: %s\n", sel)
            continue
        }
        path := entries[idx-1].Path
        err = os.RemoveAll(path)
        if err != nil {
            fmt.Printf("❌ Gagal hapus: %s - %v\n", path, err)
            continue
        }
        fmt.Printf("✅ Dihapus: %s\n", path)
        logActivity(fmt.Sprintf("Dihapus cache: %s", path))
    }
}

func diskInfo() {
    fmt.Println("=====================================")
    fmt.Println("           DISK USAGE INFO")
    fmt.Println("=====================================")
    dfCmd := exec.Command("df", "-h", "/")
    dfOut, err := dfCmd.Output()
    if err != nil {
        fmt.Println("❌ Gagal menjalankan df:", err)
    } else {
        lines := strings.Split(string(dfOut), "\n")
        if len(lines) > 1 {
            fmt.Println("Root Filesystem:")
            fmt.Println("  " + strings.Join(strings.Fields(lines[1])[1:], " "))
        }
    }
    cachePath := filepath.Join(os.Getenv("HOME"), cacheDir)
    info, err := os.Stat(cachePath)
    if err == nil && info.IsDir() {
        size, err := dirSize(cachePath)
        if err == nil {
            fmt.Printf("\nCache Directory (~/.cache):\n  Total Size: %s\n", byteCountDecimal(size))
        }
    }
    fmt.Println("=====================================")
}

func showLog() {
    fmt.Println("=====================================")
    fmt.Println("           LOG AKTIVITAS")
    fmt.Println("=====================================")
    data, err := ioutil.ReadFile(logFilePath())
    if err != nil {
        fmt.Println("ℹ Belum ada aktivitas penghapusan cache.")
    } else {
        lines := strings.Split(string(data), "\n")
        start := 0
        if len(lines) > 20 {
            start = len(lines) - 20
        }
        for _, line := range lines[start:] {
            if line != "" {
                fmt.Println(line)
            }
        }
    }
    fmt.Println("=====================================")
}

func helpMenu() {
    fmt.Println("=====================================")
    fmt.Println("              BANTUAN")
    fmt.Println("=====================================")
    fmt.Println("1. List Cache - Menampilkan daftar folder cache di ~/.cache yang ukurannya di atas batas minimum.")
    fmt.Println("2. Manage Cache - Mengelola penghapusan cache dengan memilih nomor folder cache.")
    fmt.Println("3. Exit Program - Keluar dari program.")
    fmt.Println("4. Tampilkan Info Disk Usage - Menampilkan ringkasan penggunaan disk dan cache.")
    fmt.Println("5. Log Aktivitas - Melihat aktivitas penghapusan cache terakhir.")
    fmt.Println("6. Bantuan - Menampilkan panduan ini.")
    fmt.Println("=====================================")
}

// Bantu: baca direktori cache, ukur ukuran tiap folder dan filter ukuran > minSize
func getCacheEntries(root string) ([]CacheEntry, error) {
    var entries []CacheEntry
    files, err := ioutil.ReadDir(root)
    if err != nil {
        return nil, err
    }
    for _, f := range files {
        fullPath := filepath.Join(root, f.Name())
        size, err := dirSize(fullPath)
        if err != nil {
            continue
        }
        if size >= minSize {
            entries = append(entries, CacheEntry{Path: fullPath, Size: size})
        }
    }
    sort.Slice(entries, func(i, j int) bool {
        return entries[i].Size < entries[j].Size
    })
    return entries, nil
}

// Hitung ukuran folder (rekursif)
func dirSize(path string) (int64, error) {
    var size int64
    err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            info, err := d.Info()
            if err != nil {
                return err
            }
            size += info.Size()
        }
        return nil
    })
    return size, err
}

// Hitung total size
func totalSize(entries []CacheEntry) int64 {
    var total int64
    for _, e := range entries {
        total += e.Size
    }
    return total
}

// Format bytes ke string readable (KB, MB, GB)
func byteCountDecimal(b int64) string {
    const unit = 1000
    if b < unit {
        return fmt.Sprintf("%d B", b)
    }
    div, exp := int64(unit), 0
    for n := b / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// Path log file di folder eksekusi
func logFilePath() string {
    exePath, err := os.Executable()
    if err != nil {
        return logFile
    }
    dir := filepath.Dir(exePath)
    return filepath.Join(dir, logFile)
}

func logActivity(msg string) {
    f, err := os.OpenFile(logFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return
    }
    defer f.Close()
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    f.WriteString(fmt.Sprintf("%s - %s\n", timestamp, msg))
}
