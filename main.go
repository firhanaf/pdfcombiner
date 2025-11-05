package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== PDF Combiner CLI ===")
	fmt.Println("Ketik 'all' untuk menggabungkan semua file PDF di folder saat ini.")
	fmt.Println("Atau masukkan nama file PDF satu per satu, lalu ketik 'done' jika selesai.\n")

	var inputFiles []string

	// Langkah 1: Input file
	fmt.Print("Masukkan pilihan (nama file / 'all'): ")
	firstInput, _ := reader.ReadString('\n')
	firstInput = strings.TrimSpace(firstInput)

	if strings.ToLower(firstInput) == "all" {
		// Ambil semua file PDF di folder kerja
		files, err := filepath.Glob("*.pdf")
		if err != nil {
			log.Fatalf("Gagal membaca folder: %v\n", err)
		}

		if len(files) == 0 {
			fmt.Println("⚠️  Tidak ada file PDF di folder ini.")
			return
		}

		inputFiles = append(inputFiles, files...)
		fmt.Println("📂 File PDF yang ditemukan:")
		for _, f := range files {
			fmt.Println(" -", f)
		}
	} else {
		// Tambahkan file pertama jika bukan 'all'
		if _, err := os.Stat(firstInput); os.IsNotExist(err) {
			fmt.Printf("⚠️  File '%s' tidak ditemukan.\n", firstInput)
		} else {
			inputFiles = append(inputFiles, firstInput)
		}

		// Minta input tambahan
		for {
			fmt.Print("Masukkan nama file PDF berikutnya (ketik 'done' jika selesai): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if strings.ToLower(input) == "done" {
				break
			}

			if _, err := os.Stat(input); os.IsNotExist(err) {
				fmt.Printf("⚠️  File '%s' tidak ditemukan, silakan ulangi.\n", input)
				continue
			}

			inputFiles = append(inputFiles, input)
		}
	}

	if len(inputFiles) < 2 {
		fmt.Println("❌ Minimal dua file PDF dibutuhkan untuk digabungkan.")
		return
	}

	// Langkah 2: Input nama file output
	fmt.Print("\nMasukkan nama file output (contoh: hasil.pdf): ")
	output, _ := reader.ReadString('\n')
	output = strings.TrimSpace(output)

	if !strings.HasSuffix(output, ".pdf") {
		output += ".pdf"
	}

	// Langkah 3: Gabungkan PDF
	fmt.Println("\n🔄 Menggabungkan file PDF...")
	for _, f := range inputFiles {
		fmt.Printf(" - %s\n", f)
	}

	// Gunakan konfigurasi relaxed agar tidak error metadata invalid
	conf := model.NewDefaultConfiguration()
	conf.ValidationMode = model.ValidationRelaxed

	err := api.MergeCreateFile(inputFiles, output, false, conf)
	if err != nil {
		log.Fatalf("Gagal menggabungkan PDF: %v\n", err)
	}

	fmt.Printf("\n✅ PDF berhasil digabungkan menjadi: %s\n", output)
}
