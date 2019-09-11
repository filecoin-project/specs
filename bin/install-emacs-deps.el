#!/usr/bin/env emacs -Q --script

;; Sandbox
(setq
 user-emacs-directory (concat (file-name-directory load-file-name) ".emacs/")
 package-user-dir (concat user-emacs-directory "elpa/")
 use-package-always-ensure t)

(require 'package)

; (package-initialize)

;; enable melpa
(add-to-list
 'package-archives
 '("melpa" . "https://melpa.org/packages/")
 t)

;; Update package list
(package-initialize)
(package-refresh-contents)

;; Install packages
(package-install 'use-package)
(use-package ox-hugo)
