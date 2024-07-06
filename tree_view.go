package simpleconsoleui

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zecarneiro/golangutils"
)

var (
	root         *tview.TreeNode
	tree         *tview.TreeView
	flexTree     *tview.Flex
	treeSelected string
)

// A helper function which adds the files and directories of the given path
// to the given target node.
func add(target *tview.TreeNode, path string, showDirOnly bool) {
	if showDirOnly {
		files, err := golangutils.ReadDir(path)
		golangutils.ProcessError(err)
		for _, dir := range files.Directories {
			node := tview.NewTreeNode(dir).SetReference(filepath.Join(path, dir)).SetSelectable(true)
			target.AddChild(node)
		}
	} else {
		files, err := os.ReadDir(path)
		golangutils.ProcessError(err)
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).SetReference(filepath.Join(path, file.Name())).SetSelectable(true)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}
}

func processRootDir(rootDir string, showDirOnly bool) {
	root = tview.NewTreeNode(rootDir).SetColor(tcell.ColorRed)
	tree.SetRoot(root).SetCurrentNode(root)

	// Add the current directory to the root node.
	add(root, rootDir, showDirOnly)
	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		treeSelected = reference.(string)
		children := node.GetChildren()
		if len(children) == 0 {
			file, err := os.Open(treeSelected)
			golangutils.ProcessError(err)
			fileInfo, err := file.Stat()
			golangutils.ProcessError(err)
			// Load and show files in this directory.
			if fileInfo.IsDir() {
				add(node, treeSelected, showDirOnly)
			}
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	if flexTree.GetItemCount() > 1 {
		formItem := flexTree.GetItem(0)
		flexTree.Clear()
		flexTree.AddItem(formItem, 0, 1, false)
	}
	flexTree.AddItem(tree, 0, 1, false)
}

func SelectTreeView(rootDir string, showDirOnly bool, showBorder bool, borderTitle string, callback func(selected string)) tview.Primitive {
	var dropdownList []string
	tree = tview.NewTreeView()
	flexTree = tview.NewFlex()
	formTree := tview.NewForm()
	if golangutils.IsWindows() {
		dropdownList = golangutils.GetDrives()
	} else {
		dropdownList = []string{"/"}
	}
	if len(rootDir) > 0 {
		dropdownList = append(dropdownList, rootDir)
	}

	formTree.AddDropDown("Select an root dir (hit Enter)", dropdownList, -1, func(option string, optionIndex int) {
		if optionIndex >= 0 {
			processRootDir(option, showDirOnly)
		}
	})
	formTree.AddButton("Save", func() {
		formTree.GetButton(0).Blur()
		callback(treeSelected)
	})
	formTree.SetBorder(showBorder)
	if len(borderTitle) > 0 {
		formTree.SetTitle("Select Root dir and save selected").SetTitleAlign(tview.AlignLeft)
	}
	flexTree.AddItem(formTree, 0, 1, false)
	return flexTree
}
