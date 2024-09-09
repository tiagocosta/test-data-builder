/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tiagocosta/test-data-builder/internal/builder"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "test-data-builder",
	Short: "This tool helps developers create test data in an easy, flexible and simple way",
	Long: `This tool scans your code looking for structs in order to automatically create builders for them.
The created builders are under the testdatabuilder folder.
When you use the builders, you have the option of creating your structs from them with 
the fields filled with random values ​​or default values or even to define wich value each field must have.
For example:

	Let's say that you have the following struct in your code:
	
	type User struct {
		Name string
		Age int
	} 
	
	An UserBuilder will be generated for you with the following options of usage:

	1) Generating an empty struct:
		userBulder := NewUserBuilder()
		user := userBuilder.Empty()

	2) Generating a struct with fields filled with random values:
		userBulder := NewUserBuilder()
		user := userBuilder.Build()

	3) If necessary, you can change the value of any field according to your needs, keeping the others with the generated random values:
		userBulder := NewUserBuilder()
		user := userBuilder.
			WithName("Tiago").
			Build()
	`,
	Run: func(cmd *cobra.Command, args []string) {
		gen := builder.NewGenerator()
		gen.Generate()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
