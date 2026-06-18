// Handlers para ejecucion de comandos y peticiones HTTP.
package actions

import (
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// HandlerExec ejecuta un comando y captura stdout.
// Formato: "exec:ls -la" o "exec:/usr/bin/curl https://..."
// El output se escribe en el elemento con id="exec-log" si existe.
func HandlerExec(args string, cb *Callbacks) error {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		return nil
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil && output == "" {
		output = err.Error()
	}
	if cb.SetElementText != nil {
		cb.SetElementText("exec-log", "$ "+args+"\n"+output+"\n")
	}
	cb.Log("exec: %s -> %d bytes", args, len(output))
	if cb.PublishEvent != nil {
		cb.PublishEvent("exec-output", "$ "+args+"\n"+output)
	}
	return nil
}

// HandlerCurl hace una peticion HTTP GET.
// Formato: "curl:https://api.example.com"
func HandlerCurl(args string, cb *Callbacks) error {
	url := strings.TrimSpace(args)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	output := string(body)
	if cb.SetElementText != nil {
		cb.SetElementText("exec-log", "$ curl "+args+"\n"+output+"\n")
	}
	cb.Log("curl: %s -> %d bytes", args, len(output))
	if cb.PublishEvent != nil {
		cb.PublishEvent("exec-output", "$ curl "+args+"\n"+output)
	}
	return nil
}
