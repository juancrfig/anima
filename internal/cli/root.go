package cli

import (
    "os"
    "github.com/spf13/cobra"
)


const animaArt = `
                                                      
                     @@@@@@@@@@@@@@                   
                 @@@@@@@@@@@@@@@@@@@@@                
               @@@@@@@@@@@@@@@@@@@@@@@@@@             
             @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@           
           @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@          
          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@         
         @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@        
         @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@       
        @@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@       
        @@@@@@@@@@@@@@   @@                
        @@@@@@@@@@@@       @@@    @@@@@@@@@@@@@       
                             @@ @@@@@@@@@@@@@@@       
         @@@@@@@@@@@@@@@    @@@@@@@@@@@@@@@@@@@       
         @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@        
          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@         
           @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@          
             @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@           
               @@@@@@@@@@@@@@@@@@@@@@@@@@             
                 @@@@@@@@@@@@@@@@@@@@@                
                     @@@@@@@@@@@@@@                   
`

func RootCmd() *cobra.Command {
    cmd := &cobra.Command {
            Use: "anima [command] [flags]",
            Long: animaArt + `
    This is a command-line tool that serves you as a simple personal journal. 
    You can write your diary entries, and they will be saved securely in a JSON file on your local device.
    The more you write, the better you and Anima will get to know yourself.`,
            SilenceUsage: true,
    }

    cmd.CompletionOptions.DisableDefaultCmd = true

    cmd.AddCommand(ConfigCmd())
    cmd.AddCommand(TodayCmd())

    return cmd
}

func Execute() {
    if err := RootCmd().Execute(); err != nil {
        os.Exit(1)
    }
}
