#!/usr/bin/env -S emacs -Q --script

;; load info
(load (concat (file-name-directory load-file-name) "emacs-init-build.el"))

;; process input
(with-temp-buffer
  (progn
    (condition-case nil
    (let (line)
      (while (setq line (read-from-minibuffer ""))
        (insert line)
        (insert "\n")))
      (error nil))
    (princ (org-export-as 'hugo))))
