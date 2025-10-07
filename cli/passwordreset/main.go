package passwordreset

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"

	"github.com/zobtube/zobtube/internal/config"
	"github.com/zobtube/zobtube/internal/model"
)

const userIdFromCmd = "user-id"

type Parameters struct{}

func Run(cmd *cli.Command, logger *zerolog.Logger) error {
	// setup configuration
	logger.Log().Msg("setting up configuration")
	cfg, err := config.New(
		logger,
		cmd.String("server-bind"),
		cmd.String("db-driver"),
		cmd.String("db-connstring"),
		cmd.String("media-path"),
	)
	if err != nil {
		return err
	}

	// init database
	logger.Log().Msg("initializing database connection")
	db, err := model.New(cfg)
	if err != nil {
		return err
	}

	// get user id from prompt
	userID := cmd.String(userIdFromCmd)

	// if user id is empty, display the full list of user with their ids
	if userID == "" {
		logger.Log().Msg("get user list")
		users := []model.User{}
		result := db.Find(&users)
		if result.RowsAffected < 1 {
			logger.Warn().Msg("no user found. Check your database configuration")
			return nil
		}

		for _, user := range users {
			logger.Log().Msg(fmt.Sprintf("* ID: %s (username: %s)", user.ID, user.Username))
		}

		logger.Log().Msg("please now use the --user-id flag to select the user")
		return nil
	}

	logger.Log().Str(userIdFromCmd, userID).Msg("get selected user")
	user := &model.User{
		ID: userID,
	}
	result := db.First(user)
	if result.RowsAffected < 1 {
		logger.Log().Msg("user not found, check the user id")
		return nil
	}

	newPassword := generatePassword()
	logger.Log().Str(userIdFromCmd, userID).Msg(fmt.Sprintf("new password for user %s will be %s", user.Username, newPassword))

	passwordHex := sha256.Sum256([]byte(newPassword))
	password := hex.EncodeToString(passwordHex[:])

	user.Password = password

	err = db.Save(&user).Error
	if result.RowsAffected < 1 {
		logger.Error().Err(err).Msg("unable to save new password")
		return err
	}

	logger.Info().Str(userIdFromCmd, userID).Msg("new password set successfully")

	return nil
}

func generatePassword() string {
	return rand.Text()
}
