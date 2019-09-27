;; Sandbox
(setq
 user-emacs-directory (concat (file-name-directory load-file-name) ".emacs/")
 package-user-dir (concat user-emacs-directory "elpa/")
 use-package-always-ensure t
 inhibit-message t) ; if there are errors, remove this.
 ; debug-on-error t) ; if there are errors, add this.

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
  (load (expand-file-name "deps/quicklisp/slime-helper.el"))
  (setq inferior-lisp-program "bin/sbclw")
  (add-to-list 'slime-contribs 'slime-repl))

(org-babel-do-load-languages
 'org-babel-load-languages
 '((lisp . t)))

