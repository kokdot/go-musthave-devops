package store

import (
	"bufio"
	"encoding/json"
	"os"
)

//------------------------------------producer--------------------------------------

type producer struct {
    file *os.File
    // добавляем writer в Producer
    writer *bufio.Writer
}
func NewProducer(filename string) (*producer, error) {
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0777)
    if err != nil {
        return nil, err
    }

    return &producer{
        file: file,
        // создаём новый Writer
        writer: bufio.NewWriter(file),
    }, nil
}

func (p *producer) WriteStorage(storeMap *StoreMap) error {
	p.file.Truncate(0)
    data, err := json.Marshal(storeMap)
    if err != nil {
        return err
    }
    // записываем событие в буфер
    if _, err := p.writer.Write(data); err != nil {
        return err
    }

    // добавляем перенос строки
    if err := p.writer.WriteByte('\n'); err != nil {
        return err
    }

    // записываем буфер в файл

    
    
    // return 
    p.writer.Flush()
    return p.file.Close()
}
//------------------------------------consumer--------------------------------------

type consumer struct {
    file *os.File
    // заменяем reader на scanner
    scanner *bufio.Scanner
}

func NewConsumer(filename string) (*consumer, error) {
    file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
    if err != nil {
        return nil, err
    }

    return &consumer{
        file: file,
        // создаём новый scanner
        scanner: bufio.NewScanner(file),
    }, nil
}
func (c *consumer) ReadStorage() (*StoreMap, error) {
    // одиночное сканирование до следующей строки
    if !c.scanner.Scan() {
        return nil, c.scanner.Err()
    }
    // читаем данные из scanner
    data := c.scanner.Bytes()
    sm := StoreMap{}
    err := json.Unmarshal(data, &sm)
    if err != nil {
        return nil, err
    }
    c.file.Close()
    return &sm, nil
}


func (c *consumer) Close() error {
    return c.file.Close()
}
