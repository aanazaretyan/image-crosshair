package processor

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// LineWidth определяет ширину линий в пикселях
const LineWidth = 2

// ProcessImage обрабатывает одно изображение, рисуя на нем вертикальную и горизонтальную линии
func ProcessImage(inputPath, outputPath string) error {
	// Открываем файл изображения
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	// Декодируем изображение
	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("не удалось декодировать изображение: %w", err)
	}

	// Создаем новое изображение RGBA для рисования
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	rgba := image.NewRGBA(bounds)

	// Копируем исходное изображение в RGBA
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Рисуем вертикальную линию посередине
	centerX := width / 2
	for x := centerX - LineWidth/2; x < centerX+LineWidth/2; x++ {
		for y := 0; y < height; y++ {
			rgba.Set(x, y, color.Black)
		}
	}

	// Рисуем горизонтальную линию посередине
	centerY := height / 2
	for y := centerY - LineWidth/2; y < centerY+LineWidth/2; y++ {
		for x := 0; x < width; x++ {
			rgba.Set(x, y, color.Black)
		}
	}

	// Создаем выходной файл
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("не удалось создать выходной файл: %w", err)
	}
	defer outFile.Close()

	// Сохраняем изображение в том же формате, что и исходное
	switch strings.ToLower(format) {
	case "jpeg":
		err = jpeg.Encode(outFile, rgba, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(outFile, rgba)
	default:
		// Если формат не поддерживается, сохраняем как PNG
		err = png.Encode(outFile, rgba)
	}

	if err != nil {
		return fmt.Errorf("не удалось сохранить изображение: %w", err)
	}

	return nil
}

// ProcessDirectory обрабатывает все изображения в указанной директории
func ProcessDirectory(inputDir, outputDir string) error {
	// Получаем список файлов в директории
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return fmt.Errorf("не удалось прочитать директорию: %w", err)
	}

	// Обрабатываем каждый файл
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Пропускаем поддиректории
		}

		// Проверяем, является ли файл изображением
		filename := entry.Name()
		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue // Пропускаем не-изображения
		}

		// Формируем пути к файлам
		inputPath := filepath.Join(inputDir, filename)
		outputPath := filepath.Join(outputDir, filename)

		// Обрабатываем изображение
		if err := ProcessImage(inputPath, outputPath); err != nil {
			fmt.Printf("Ошибка при обработке %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Обработан файл: %s\n", filename)
	}

	return nil
}

// IsImageFile проверяет, является ли файл изображением по расширению
func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
