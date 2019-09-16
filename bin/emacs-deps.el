;; require package
(require 'package)

;; enable melpa, if not there.
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

;; slime is used in orient + orgmode
(use-package slime
  :init
  (global-set-key (kbd "C-c z") 'slime-repl)
  (load (expand-file-name ".quicklisp/slime-helper.el"))
  (setq inferior-lisp-program "bin/sbcl")
  (add-to-list 'slime-contribs 'slime-repl))
