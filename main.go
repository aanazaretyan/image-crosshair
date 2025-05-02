package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aanazaretyan/image-crosshair/internal/processor"
)

func main() {
	// Парсим аргументы командной строки
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Ошибка: не указаны входные файлы или директория")
		fmt.Println("Использование: image-crosshair [файл1 файл2 ...] или [директория]")
		os.Exit(1)
	}

	// Создаем директорию для результатов с текущей датой и временем
	timestamp := time.Now().Format("20060102-15:04:05")
	outputDir := fmt.Sprintf("result-%s", timestamp)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Ошибка при создании директории результатов: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Результаты будут сохранены в директорию: %s\n", outputDir)

	// Проверяем, является ли первый аргумент директорией
	fileInfo, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("Ошибка при получении информации о файле/директории: %v\n", err)
		os.Exit(1)
	}

	if fileInfo.IsDir() {
		// Обрабатываем все изображения в директории
		fmt.Printf("Обработка изображений в директории: %s\n", args[0])
		if err := processor.ProcessDirectory(args[0], outputDir); err != nil {
			fmt.Printf("Ошибка при обработке директории: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Обрабатываем список файлов
		for _, inputPath := range args {
			if !processor.IsImageFile(inputPath) {
				fmt.Printf("Пропуск файла (не изображение): %s\n", inputPath)
				continue
			}

			filename := filepath.Base(inputPath)
			outputPath := filepath.Join(outputDir, filename)

			fmt.Printf("Обработка файла: %s\n", inputPath)
			if err := processor.ProcessImage(inputPath, outputPath); err != nil {
				fmt.Printf("Ошибка при обработке %s: %v\n", inputPath, err)
				continue
			}
		}
	}

	fmt.Println("Обработка завершена. Результаты сохранены в директорию:", outputDir)
}
