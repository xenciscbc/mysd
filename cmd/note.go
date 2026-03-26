package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xenciscbc/mysd/internal/spec"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage deferred notes",
	RunE:  runNoteList,
}

var noteAddCmd = &cobra.Command{
	Use:   "add [content...]",
	Short: "Add a deferred note",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runNoteAdd,
}

var noteDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a deferred note by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runNoteDelete,
}

func init() {
	rootCmd.AddCommand(noteCmd)
	noteCmd.AddCommand(noteAddCmd)
	noteCmd.AddCommand(noteDeleteCmd)
}

func runNoteList(cmd *cobra.Command, args []string) error {
	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	store, err := spec.LoadDeferredStore(specDir)
	if err != nil {
		return err
	}

	if len(store.Notes) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No deferred notes.")
		return nil
	}

	for _, note := range store.Notes {
		fmt.Fprintf(cmd.OutOrStdout(), "[%d] %s  (%s)\n", note.ID, note.Content, note.CreatedAt)
	}
	return nil
}

func runNoteAdd(cmd *cobra.Command, args []string) error {
	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	content := strings.Join(args, " ")

	store, err := spec.LoadDeferredStore(specDir)
	if err != nil {
		return err
	}

	note := store.Add(content)

	if err := spec.SaveDeferredStore(specDir, store); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Added note #%d: %s\n", note.ID, note.Content)
	return nil
}

func runNoteDelete(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid ID %q: must be a number", args[0])
	}

	specDir, _, err := spec.DetectSpecDir(".")
	if err != nil {
		return err
	}

	store, err := spec.LoadDeferredStore(specDir)
	if err != nil {
		return err
	}

	if !store.Delete(id) {
		return fmt.Errorf("note #%d not found", id)
	}

	if err := spec.SaveDeferredStore(specDir, store); err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Deleted note #%d\n", id)
	return nil
}
