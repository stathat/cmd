// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

const logo1 = `
/ _\ |_ __ _| |_  /\  /\__ _| |_
\ \| __/ _' | __|/ /_/ / _' | __|
_\ \ || (_| | |_/ __  / (_| | |_ 
\__/\__\__,_|\__\/ /_/ \__,_|\__|
`
const logo2 = `
 #####                     #     #              
#     # #####   ##   ##### #     #   ##   ##### 
#         #    #  #    #   #     #  #  #    #   
 #####    #   #    #   #   ####### #    #   #   
      #   #   ######   #   #     # ######   #   
#     #   #   #    #   #   #     # #    #   #   
 #####    #   #    #   #   #     # #    #   #   
`

const logo3 = `
     _____                               _____                    
  __|___  |__    __    ____     __    __|  _  |__  ____     __    
 |   ___|    | _|  |_ |    \  _|  |_ |  |_| |    ||    \  _|  |_  
  '-.'-.     ||_    _||     \|_    _||   _  |    ||     \|_    _| 
 |______|  __|  |__|  |__|\__\ |__|  |__| |_|  __||__|\__\ |__|   
    |_____|                             |_____|                   
`

const logo4 = `
   .dMMMb dMMMMMMP .aMMMb dMMMMMMP dMP dMP .aMMMb dMMMMMMP 
  dMP" VP   dMP   dMP"dMP   dMP   dMP dMP dMP"dMP   dMP    
  VMMMb    dMP   dMMMMMP   dMP   dMMMMMP dMMMMMP   dMP     
dP .dMP   dMP   dMP dMP   dMP   dMP dMP dMP dMP   dMP      
VMMMP"   dMP   dMP dMP   dMP   dMP dMP dMP dMP   dMP       
`

const logo5 = `
 __                  
(_ _|_ _ _|_|_| _ _|_
__) |_(_| |_| |(_| |_
`

const logo6 = `
    ___       ___       ___       ___       ___       ___       ___   
   /\  \     /\  \     /\  \     /\  \     /\__\     /\  \     /\  \  
  /::\  \    \:\  \   /::\  \    \:\  \   /:/__/_   /::\  \    \:\  \ 
 /\:\:\__\   /::\__\ /::\:\__\   /::\__\ /::\/\__\ /::\:\__\   /::\__\
 \:\:\/__/  /:/\/__/ \/\::/  /  /:/\/__/ \/\::/  / \/\::/  /  /:/\/__/
  \::/  /   \/__/      /:/  /   \/__/      /:/  /    /:/  /   \/__/   
   \/__/               \/__/               \/__/     \/__/            
`

const logo7 = `
                                                                          
      _/_/_/    _/                  _/      _/    _/              _/      
   _/        _/_/_/_/    _/_/_/  _/_/_/_/  _/    _/    _/_/_/  _/_/_/_/   
    _/_/      _/      _/    _/    _/      _/_/_/_/  _/    _/    _/        
       _/    _/      _/    _/    _/      _/    _/  _/    _/    _/         
_/_/_/        _/_/    _/_/_/      _/_/  _/    _/    _/_/_/      _/_/      
`

const logo8 = `
_       ___        ______  ______        __    ___    _____  ______        __
 )  ____) (__    __)    /  \    (__    __) \  |   |  /    /  \    (__    __) 
(  (___      |  |      /    \      |  |     |  \_/  |    /    \      |  |    
 \___  \     |  |     /  ()  \     |  |     |   _   |   /  ()  \     |  |    
 ____)  )    |  |    |   __   |    |  |     |  / \  |  |   __   |    |  |    
(      (_____|  |____|  (__)  |____|  |____/  |___|  \_|  (__)  |____|  |____
`

const logo9 = `
 __           o        
(_ _|_ _ _|_  |  |_| _ _|_
__) |_(_| |_  |  | |(_| |_
              o
`

var logos = []string{logo1, logo2, logo3, logo4, logo5, logo6, logo7, logo8, logo9}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print stathat cmd version",
	Long:  `Version prints the stathat command line utility version`,
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(time.Now().Unix())
		index := rand.Intn(len(logos))
		fmt.Println(logos[index])
		fmt.Printf("stathat cmd version 0.1.0\n\n")
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
