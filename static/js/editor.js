function createEditorMd(divId, submitId, markdown) {
    var editor = editormd(divId, {
        height: 400,
		markdown: markdown,
	    autoFocus: false,
        path: "/static/js/editor.md-1.5.0/lib/",
	    placeholder: "采用markdown语法",
        toolbarIcons: function() {
          return ["undo", "redo", "|", "bold", "italic", "quote", "|", "h1", "h2", "h3", "h4", "h5", "h6", "|", "list-ul", "list-ol", "hr", "|", "link", "reference-link", "image", "code", "preformatted-text", "code-block", "|", "goto-line", "watch", "preview", "fullscreen", "|", "help", "info"]
        },
        saveHTMLToTextarea: true,
        imageUpload: false,
        //imageFormats: ["jpg", "jpeg", "gif", "png"],
        //imageUploadURL: "/upload/image",
	    onchange: function() {
	      $(submitId).attr('disabled', this.getMarkdown().trim() == "");
	    }
      });

	return editor;
}