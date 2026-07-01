package mrutils

import (
	"fmt"
	"io"

	gitlab "gitlab.com/gitlab-org/api/client-go/v2"

	"gitlab.com/gitlab-org/cli/internal/iostreams"
	"gitlab.com/gitlab-org/cli/internal/utils"
)

// noteUsername returns the note author's username, falling back to "unknown"
// if the username is empty (e.g. redacted users or system-generated notes).
func noteUsername(n *gitlab.Note) string {
	if n.Author.Username != "" {
		return n.Author.Username
	}
	return "unknown"
}

// noteTimeAgo returns a human-readable "time ago" string for the note's
// creation time, or an empty string if CreatedAt is nil.
func noteTimeAgo(n *gitlab.Note) string {
	if n.CreatedAt != nil {
		return utils.TimeToPrettyTimeAgo(*n.CreatedAt)
	}
	return ""
}

// renderBody renders the body as markdown for a terminal, or returns it
// verbatim when piped so the raw output stays free of control characters.
func renderBody(ios *iostreams.IOStreams, body string) string {
	if !ios.IsOutputTTY() {
		return body
	}
	rendered, err := utils.RenderMarkdown(body, ios.BackgroundColor())
	if err != nil {
		return body
	}
	return rendered
}

// PrintDiscussions renders discussions to out.
func PrintDiscussions(out io.Writer, ios *iostreams.IOStreams, discussions []*gitlab.Discussion, showSystemLogs bool) {
	c := ios.Color()

	for _, discussion := range discussions {
		if len(discussion.Notes) == 0 {
			continue
		}

		firstNote := discussion.Notes[0]

		// Skip system notes unless showSystemLogs is set
		if firstNote.System && !showSystemLogs {
			continue
		}

		// Threaded discussions (not individual notes)
		if !discussion.IndividualNote && len(discussion.Notes) > 1 {
			fmt.Fprintf(out, "Thread [discussion: %s]", TruncateDiscussionID(discussion.ID))

			// Show resolution status if resolvable
			if firstNote.Resolvable {
				if firstNote.Resolved {
					fmt.Fprint(out, c.Green(" ✓ resolved"))
				} else {
					fmt.Fprint(out, c.Yellow(" ⚠ unresolved"))
				}
			}
			fmt.Fprintln(out)

			// Print first note
			createdAt := noteTimeAgo(firstNote)
			fmt.Fprintf(out, "  @%s commented ", noteUsername(firstNote))
			fmt.Fprintf(out, "%s %s\n", c.Gray(createdAt), c.Gray(fmt.Sprintf("[note #%d]", firstNote.ID)))

			if firstNote.Position != nil {
				PrintCommentFileContext(out, c, firstNote.Position)
			}

			body := renderBody(ios, firstNote.Body)
			fmt.Fprintln(out, utils.Indent(body, "  "))
			fmt.Fprintln(out)

			// Print replies (indented)
			for i, note := range discussion.Notes[1:] {
				if note.System && !showSystemLogs {
					continue
				}
				replyTime := noteTimeAgo(note)
				fmt.Fprintf(out, "    @%s replied ", noteUsername(note))
				fmt.Fprintf(out, "%s %s\n", c.Gray(replyTime), c.Gray(fmt.Sprintf("[note #%d]", note.ID)))

				replyBody := renderBody(ios, note.Body)
				fmt.Fprintln(out, utils.Indent(replyBody, "    "))
				if i < len(discussion.Notes[1:])-1 {
					fmt.Fprintln(out)
				}
			}
			fmt.Fprintln(out)
		} else {
			// Individual note (not a thread)
			note := firstNote
			createdAt := noteTimeAgo(note)
			fmt.Fprint(out, "@", noteUsername(note))
			if note.System {
				fmt.Fprintf(out, " %s ", note.Body)
				fmt.Fprintln(out, c.Gray(createdAt))
			} else {
				body := renderBody(ios, note.Body)
				fmt.Fprint(out, " commented ")
				fmt.Fprintf(out, "%s %s\n", c.Gray(createdAt), c.Gray(fmt.Sprintf("[note #%d]", note.ID)))

				if note.Position != nil {
					PrintCommentFileContext(out, c, note.Position)
				}

				fmt.Fprintln(out, utils.Indent(body, " "))
			}
			fmt.Fprintln(out)
		}
	}
}

// PrintCommentFileContext prints file and line context for a note position.
func PrintCommentFileContext(out io.Writer, c *iostreams.ColorPalette, pos *gitlab.NotePosition) {
	// Check for multi-line comment first
	if pos.LineRange != nil && pos.LineRange.StartRange != nil && pos.LineRange.EndRange != nil {
		startLine := pos.LineRange.StartRange.NewLine
		endLine := pos.LineRange.EndRange.NewLine

		// Fall back to old line numbers if new ones aren't available
		if startLine == 0 {
			startLine = pos.LineRange.StartRange.OldLine
		}
		if endLine == 0 {
			endLine = pos.LineRange.EndRange.OldLine
		}

		// Display range if we have valid start and end lines
		if startLine > 0 && endLine > 0 {
			filePath := pos.NewPath
			if filePath == "" {
				filePath = pos.OldPath
			}
			if filePath != "" {
				if startLine != endLine {
					fmt.Fprintf(out, " on %s:%d-%d\n", c.Cyan(filePath), startLine, endLine)
				} else {
					fmt.Fprintf(out, " on %s:%d\n", c.Cyan(filePath), startLine)
				}
				return
			}
		}
	}

	// Fall back to single-line comment
	if pos.NewPath != "" && pos.NewLine > 0 {
		fmt.Fprintf(out, " on %s:%d\n", c.Cyan(pos.NewPath), pos.NewLine)
	} else if pos.OldPath != "" && pos.OldLine > 0 {
		fmt.Fprintf(out, " on %s:%d\n", c.Cyan(pos.OldPath), pos.OldLine)
	}
}
