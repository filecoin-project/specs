; sbcl-userinit
;
; this file is here because it's tricky to edit the userinit file
; reliably from within sbcl as called by a shell script. i decided
; to do this manually.
;
; use load-quicklisp-verbose if you need to debug quicklisp loading

(defvar *quicklisp-init-file*
  (merge-pathnames "deps/quicklisp/setup.lisp" (user-homedir-pathname)))

(defun load-quicklisp-verbose ()
  (write-line (format nil "looking for quicklisp in ~A" *quicklisp-init-file*))
  (if (probe-file *quicklisp-init-file*)
    (progn
      (write-line "quicklisp found")
      (if (load *quicklisp-init-file*)
        (write-line "quicklisp loaded")
        (write-line "Warning: quicklisp failed to load"))
      t)
    (progn
      (write-line "Warning: quicklisp not found")
      t)))

(defun load-quicklisp-silent ()
  (when (probe-file *quicklisp-init-file*)
      (load *quicklisp-init-file*)))

; actually run load-quicklisp function. comment one of these out.
; (load-quicklisp-verbose)
(load-quicklisp-silent)
