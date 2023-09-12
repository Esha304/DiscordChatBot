package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	dg *discordgo.Session
)

func main() {
	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a new Discord session
	var err error
	dg, err = discordgo.New("Bot MTEyNzMwMjkxNzcxNjcwOTM5Ng.G6pVn8.m_XqVxk8ZSMxxVAvZcRfUkF_Ddfq9pqssUnAKo")
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Register the onReady function as the event handler for the Ready event
	dg.AddHandler(onReady)

	// Register the messageCreate function as the event handler for the MessageCreate event
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}

	// Keep the bot running until interrupted
	<-make(chan struct{})
}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message starts with the command prefix "!"
	if strings.HasPrefix(m.Content, "!") {
		// Extract the command and arguments from the message content
		parts := strings.Split(m.Content[1:], " ")
		command := parts[0]
		args := parts[1:]

		// Handle the "vocabulary" command
		if command == "vocabulary" {
			handleVocabularyCommand(s, m.ChannelID, args)
		} else if command == "number" {
			handleNumberCommand(s, m.ChannelID, args)
			return
		} else if command == "exercise" {
			handleExerciseCommand(s, m.ChannelID, args)
			return
		} else {
			// Command not recognized, send an error message
			s.ChannelMessageSend(m.ChannelID, "Command not recognized.")
		}
	}
}

func handleVocabularyCommand(s *discordgo.Session, channelID string, args []string) {
	// Generate new vocabulary words
	vocabularyWords := generateVocabularyWords()

	// Send the vocabulary words as a message in the channel
	_, err := s.ChannelMessageSend(channelID, "Vocabulary Words:\n"+strings.Join(vocabularyWords, "\n"))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func generateVocabularyWords() []string {
	// Generate new vocabulary words
	vocabularyWords := []string{
		"Hola - Hello",
		"¿cómo estás? - how are you?",
		"buenos días - good morning",
		"¿entiende(s)? - do you understand?",
		"muy bien - very well",
		"no (lo) sé - i don’t know",
		"no hablo español - i don’t speak spanish",
		"Adiós - Goodbye",
	}

	// Shuffle the vocabulary words using Fisher-Yates algorithm
	for i := len(vocabularyWords) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		vocabularyWords[i], vocabularyWords[j] = vocabularyWords[j], vocabularyWords[i]
	}

	return vocabularyWords
}

func handleExerciseCommand(s *discordgo.Session, channelID string, args []string) {
	// Generate an interactive language exercise
	question, answer := generateExercise()

	// Send the exercise question as a message in the channel
	_, err := s.ChannelMessageSend(channelID, "Language Exercise:\n"+question)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}

	// Variable to track if the event handler should be active or not
	activeHandler := true

	// Create a function to handle the user's response
	handleResponse := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages sent by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Check if the message is in the same channel
		if m.ChannelID == channelID && activeHandler {
			// Compare the user's response to the answer
			isCorrect := checkAnswer(m.Content, answer)

			// Send feedback based on the correctness of the answer
			if isCorrect {
				s.MessageReactionAdd(channelID, m.ID, "✅") // Add a checkmark emoji as a reaction
				s.ChannelMessageSend(channelID, "Great job! Your answer is correct.")
			} else {
				s.MessageReactionAdd(channelID, m.ID, "❌") // Add a cross mark emoji as a reaction
				s.ChannelMessageSend(channelID, "Oops! Your answer is incorrect. Try again.")
			}

			// Set the activeHandler variable to false to deactivate the event handler
			activeHandler = false
		}
	}

	// Listen for the user's response
	s.AddHandler(handleResponse)
}

func generateExercise() (string, string) {
	// Define a list of exercise questions and answers
	exercises := map[string]string{
		"Translate 'Hello' to Spanish.":     "Hola",
		"Translate 'Goodbye' to Spanish.":   "Adiós",
		"Translate 'Thank you' to Spanish.": "Gracias",
		"Translate 'Yes' to Spanish.":       "Sí",
		"Translate 'No' to Spanish.":        "No",
	}

	// Select a random exercise question
	var question string
	var answer string
	for q, a := range exercises {
		question = q
		answer = a
		break
	}

	return question, answer
}

func checkAnswer(response, answer string) bool {
	// Prepare the expected answer by removing whitespace and converting to lowercase
	expectedAnswer := strings.ToLower(strings.TrimSpace(answer))

	// Prepare the user's response by removing whitespace and converting to lowercase
	userResponse := strings.ToLower(strings.TrimSpace(response))

	// Compare the expected answer to the user's response
	return userResponse == expectedAnswer
}

func handleNumberCommand(s *discordgo.Session, channelID string, args []string) {
	if len(args) < 1 {
		s.ChannelMessageSend(channelID, "Please provide a number.")
		return
	}

	// Get the number argument from the command
	numberStr := args[0]

	// Convert the number string to an integer
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		s.ChannelMessageSend(channelID, "Invalid number provided.")
		return
	}

	// Translate the number into Spanish
	spanishTranslation := translateNumberToSpanish(number)

	// Send the translated word as a message in the channel
	_, err = s.ChannelMessageSend(channelID, fmt.Sprintf("Spanish Translation of %d: %s", number, spanishTranslation))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func translateNumberToSpanish(number int) string {
	// Define a map of number translations in Spanish
	numberTranslations := map[int]string{
		1:  "uno",
		2:  "dos",
		3:  "tres",
		4:  "cuatro",
		5:  "cinco",
		6:  "seis",
		7:  "siete",
		8:  "ocho",
		9:  "nueve",
		10: "diez",
		11: "once",
		12: "doce",
		13: "trece",
		14: "catorce",
		15: "quince",
		16: "dieciséis",
		17: "diecisiete",
		18: "dieciocho",
		19: "diecinueve",
		20: "veinte",
		21: "veintiuno",
		22: "veintidós",
		23: "veintitrés",
		24: "veinticuatro",
		25: "veinticinco",
		26: "veintiséis",
		27: "veintisiete",
		28: "veintiocho",
		29: "veintinueve",
		30: "treinta",
	}

	// Check if the number is in the map
	if translation, ok := numberTranslations[number]; ok {
		return translation
	}

	// Return an empty string if no translation is found
	return ""
}
