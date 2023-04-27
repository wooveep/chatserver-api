/*
 * @Author: cloudyi.li
 * @Date: 2023-04-27 11:22:13
 * @LastEditTime: 2023-04-27 11:22:17
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/openai/form_builder.go
 */
package openai

import (
	"io"
	"mime/multipart"
	"os"
)

type formBuilder interface {
	createFormFile(fieldname string, file *os.File) error
	writeField(fieldname, value string) error
	close() error
	formDataContentType() string
}

type defaultFormBuilder struct {
	writer *multipart.Writer
}

func newFormBuilder(body io.Writer) *defaultFormBuilder {
	return &defaultFormBuilder{
		writer: multipart.NewWriter(body),
	}
}

func (fb *defaultFormBuilder) createFormFile(fieldname string, file *os.File) error {
	fieldWriter, err := fb.writer.CreateFormFile(fieldname, file.Name())
	if err != nil {
		return err
	}

	_, err = io.Copy(fieldWriter, file)
	if err != nil {
		return err
	}
	return nil
}

func (fb *defaultFormBuilder) writeField(fieldname, value string) error {
	return fb.writer.WriteField(fieldname, value)
}

func (fb *defaultFormBuilder) close() error {
	return fb.writer.Close()
}

func (fb *defaultFormBuilder) formDataContentType() string {
	return fb.writer.FormDataContentType()
}
