#!/bin/bash

# =====================================
#    Cache Management Script
# =====================================

TITLE="Manage
Chace"
MIN_SIZE=1K # Batas minimum cache yang akan ditampilkan
LOG_FILE="$(dirname "$0")/cache_management.log"

# Fungsi untuk konversi ukuran menjadi byte
to_bytes() {
    local size=$1
    local unit=${size: -1}
    local val=${size%?}
    case "$unit" in
        K|k) echo $(echo "$val * 1024" | bc) ;;
        M|m) echo $(echo "$val * 1024 * 1024" | bc) ;;
        G|g) echo $(echo "$val * 1024 * 1024 * 1024" | bc) ;;
        *) echo $val ;;
    esac
}

# Fungsi untuk list cache dalam bentuk tabel
list_cache() {
    echo "====================================="
    echo "         LIST CACHE > $MIN_SIZE"
    echo "====================================="
    printf "%-4s | %-8s | %s\n" "No" "Size" "Folder"
    echo "----+----------+-------------------------"

    local count=1
    while IFS=$'\t' read -r size folder; do
        printf "%-4s | %-8s | %s\n" "$count" "$size" "$folder"
        ((count++))
    done < <(du -sh "$HOME/.cache"/* 2>/dev/null | sort -h | tr -s ' ' '\t')

    echo "----+----------+-------------------------"

    # Hitung total cache
    TOTAL_CACHE=$(du -sh "$HOME/.cache" 2>/dev/null | cut -f1)
    echo "💾 Total Cache: $TOTAL_CACHE"
    echo "====================================="
}

# Fungsi untuk manage cache (tabel + pilih nomor)
manage_cache() {
    echo "====================================="
    echo "      MANAGE CACHE > $MIN_SIZE"
    echo "====================================="

    # Ambil daftar cache di atas MIN_SIZE
    mapfile -t CACHE_LIST < <(
        du -sh "$HOME/.cache"/* 2>/dev/null | \
        sort -h | \
        awk -v min="$MIN_SIZE" '
        function to_bytes(size) {
            unit = substr(size, length(size))
            val = substr(size, 1, length(size)-1)
            if (unit == "K") return val*1024
            if (unit == "M") return val*1024*1024
            if (unit == "G") return val*1024*1024*1024
            return val
        }
        {
            if (to_bytes($1) >= to_bytes(min)) print $1 "\t" $2
        }'
    )

    if [ ${#CACHE_LIST[@]} -eq 0 ]; then
        echo "✅ Tidak ada cache di atas $MIN_SIZE."
        return
    fi

    # Tampilkan tabel
    printf "%-4s | %-8s | %s\n" "No" "Size" "Folder"
    echo "----+----------+--------------------------------"
    for i in "${!CACHE_LIST[@]}"; do
        size=$(echo "${CACHE_LIST[$i]}" | cut -f1)
        folder=$(echo "${CACHE_LIST[$i]}" | cut -f2-)
        printf "%-4s | %-8s | %s\n" "$((i+1))" "$size" "$folder"
    done
    echo "----+----------+--------------------------------"

    # Input nomor
    echo "📌 Ketik nomor cache yang ingin dihapus (pisahkan spasi)"
    echo "📌 Ketik Q untuk kembali ke menu utama"
    read -p "Masukkan pilihan: " pilihan

    # Keluar jika Q/q
    [[ "$pilihan" =~ ^[Qq]$ ]] && { echo "🔙 Kembali ke menu utama..."; return; }

    # Loop nomor yang dimasukkan
    for num in $pilihan; do
        if [[ "$num" =~ ^[0-9]+$ ]] && (( num >= 1 && num <= ${#CACHE_LIST[@]} )); then
            folder=$(echo "${CACHE_LIST[$((num-1))]}" | cut -f2-)
            rm -rf "$folder"
            echo "✅ Dihapus: $folder"
            echo "$(date '+%Y-%m-%d %H:%M:%S') - Dihapus cache: $folder" >> "$LOG_FILE"
        else
            echo "⚠ Nomor tidak valid: $num"
        fi
    done
}

# Fungsi tampilkan info disk usage
disk_info() {
    echo "====================================="
    echo "           DISK USAGE INFO"
    echo "====================================="
    echo "Root Filesystem:"
    df -h / | awk 'NR==2 {print "  Size: "$2", Used: "$3", Avail: "$4", Use%: "$5}'
    echo
    echo "Cache Directory (~/.cache):"
    du -sh "$HOME/.cache" 2>/dev/null | awk '{print "  Total Size: "$1}'
    echo "====================================="
}

# Fungsi tampilkan log aktivitas
show_log() {
    echo "====================================="
    echo "           LOG AKTIVITAS"
    echo "====================================="
    if [ ! -f "$LOG_FILE" ] || [ ! -s "$LOG_FILE" ]; then
        echo "ℹ Belum ada aktivitas penghapusan cache."
    else
        tail -n 20 "$LOG_FILE"
    fi
    echo "====================================="
    read -p "Tekan ENTER untuk kembali ke menu..."
}

# Fungsi bantuan / help
help_menu() {
    echo "====================================="
    echo "              BANTUAN"
    echo "====================================="
    echo "1. List Cache - Menampilkan daftar folder cache di ~/.cache yang ukurannya di atas batas minimum."
    echo "2. Manage Cache - Mengelola penghapusan cache dengan memilih nomor folder cache."
    echo "3. Exit Program - Keluar dari program."
    echo "4. Tampilkan Info Disk Usage - Menampilkan ringkasan penggunaan disk dan cache."
    echo "5. Log Aktivitas - Melihat aktivitas penghapusan cache terakhir."
    echo "6. Bantuan - Menampilkan panduan ini."
    echo "====================================="
    read -p "Tekan ENTER untuk kembali ke menu..."
}

# Main menu loop
while true; do
    clear
    echo -e "====================================="
    figlet -f standard "$TITLE"
    echo -e "====================================="
    echo -e "\e[0;32m             by - fmp\e[0m"
    echo "====================================="

    MENU_ITEMS="1. List Cache\n2. Manage Cache\n3. Exit Program\n4. Tampilkan Info Disk Usage\n5. Log Aktivitas\n6. Bantuan"

    pilihan=$(echo -e "$MENU_ITEMS" | fzf --height=10 --border --prompt="Pilih menu: ")

    case "$pilihan" in
        1*|List*) list_cache ;;
        2*|Manage*) manage_cache ;;
        3*|Exit*) echo "👋 Keluar dari program..."; exit 0 ;;
        4*|Tampilkan*) disk_info ;;
        5*|Log*) show_log ;;
        6*|Bantuan*) help_menu ;;
        *) echo "⚠ Pilihan tidak valid." ;;
    esac

    echo
    read -p "Tekan ENTER untuk kembali ke menu..."
done
