// Copyright (c) 2020, Arveto Ink. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package auth

import (
	"bytes"
	"encoding/base64"
	"io"
	"log"
	"mime"
	"net/http"
	"net/smtp"
)

// Send a mail to test the config.
func (s *Server) testMail(w http.ResponseWriter, r *http.Request) {
	to := r.URL.Query().Get("to")
	if to == "" {
		s.Error(w, r, "Give to in URL params", http.StatusBadRequest)
		return
	}

	if err := s.sendMail(to, "Test Mail", "A mail to test the config"); err != nil {
		s.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Send mail success!"))
}

// Send a mail.
func (s *Server) sendMail(to, subject, body string) error {
	// Header
	buff := bytes.NewBufferString("Content-Transfer-Encoding: base64\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"MIME-Version: 1.0\r\n")
	buff.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n")
	buff.WriteString("To: " + to + "\r\n\r\n")
	encodeBody(buff, []byte(body))
	buff.WriteString("\r\n")

	// Send mail
	err := smtp.SendMail(s.mailHost, s.mailAuth, s.mailLogin, []string{to}, buff.Bytes())
	if err != nil {
		log.Printf("[SEND MAIL ERROR] <%s> %v\n", to, err)
		return err
	}

	log.Printf("[MAIL] <%s> %q\n", to, subject)

	return nil
}

func encodeBody(m io.Writer, src []byte) {
	b := base64.NewEncoder(base64.StdEncoding, m)
	defer b.Close()

	const l = 74
	i := 0
	for ; i+l < len(src); i += l {
		b.Write(src[i : i+l])
		m.Write([]byte("\r\n"))
	}
	b.Write(src[i:])
}
