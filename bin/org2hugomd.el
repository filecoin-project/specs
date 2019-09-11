#!/usr/bin/env -S emacs -Q --script

;; Sandbox
(setq
 user-emacs-directory (concat (file-name-directory load-file-name) ".emacs/")
 package-user-dir (concat user-emacs-directory "elpa/")
 use-package-always-ensure t
 inhibit-message t) ; if there are errors, remove this.
 ; debug-on-error t) ; if there are errors, add this.

;; require package
(require 'package)

;; enable melpa
(add-to-list
 'package-archives
 '("melpa" . "https://melpa.org/packages/")
 t)

;; Update package list
(package-initialize)
(unless (require 'use-package nil 'noerror)
  (package-refresh-contents)
  (package-install 'use-package))

;; Load packages we need
(use-package ox-hugo)

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
