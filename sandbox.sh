MENU_ITEMS="List Cache\nManage Cache\nExit Program\nDisk Usage Info\nLog Aktivitas\nBantuan"

choice=$(echo -e "$MENU_ITEMS" | fzf --height=10 --border --prompt="Pilih menu: ")

case "$choice" in
  "List Cache") list_cache ;;
  "Manage Cache") manage_cache ;;
  "Exit Program") echo "Bye!"; exit 0 ;;
  "Disk Usage Info") disk_info ;;
  "Log Aktivitas") show_log ;;
  "Bantuan") help_menu ;;
  *) echo "Pilihan tidak valid." ;;
esac
