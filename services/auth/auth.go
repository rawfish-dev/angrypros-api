package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/rawfish-dev/angrypros-api/config"
)

var _ AuthService = new(Service)

type AuthService interface {
	CreateFirebaseUser(emailAddress, title, password string) (firebaseUserId string, err error)
	GetFirebaseUserId(idToken string) (firebaseUserId string, err error)
	GetFirebaseUserEmail(firebaseUserId string) (email string, err error)
	// VerifyRecaptcha(recaptchaToken string) (err error)
	// SendForgotPasswordEmail(email string) (err error)
}

type GoogleRecaptchaRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

type GoogleRecaptchaResponse struct {
	Success            bool     `json:"success"`
	ChallengeTimestamp string   `json:"challenge_ts"`
	Hostname           string   `json:"hostname"`
	ErrorCodes         []string `json:"error-codes"`
}

type Service struct {
	firebaseApp               *firebase.App
	skipRecaptchaVerification bool
	// recaptchaVerificationUrl  string
	// recaptchaSecret           string
}

func NewService(g config.GoogleConfig) (*Service, error) {
	googleConfigBytes, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("could not marshal google config to bytes due to %s", err)
	}

	opt := option.WithCredentialsJSON(googleConfigBytes)
	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app due to %s", err)
	}

	s := &Service{
		firebaseApp:               firebaseApp,
		skipRecaptchaVerification: true,
		// skipRecaptchaVerification: r.SkipVerification,
		// recaptchaVerificationUrl:  r.VerificationUrl,
		// recaptchaSecret:           r.Secret,
	}

	return s, nil
}

func (s Service) CreateFirebaseUser(emailAddress, username, password string) (string, error) {
	ctx := context.Background()

	authClient, err := s.firebaseApp.Auth(ctx)
	if err != nil {
		log.Printf("encountered error while creating Firebase auth client due to %s", err)
		return "", err
	}

	params := (&auth.UserToCreate{}).
		Email(emailAddress).
		DisplayName(username).
		Password(password)
	firebaseUser, err := authClient.CreateUser(ctx, params)
	if err != nil {
		log.Printf("encountered error while creating Firebase user due to %s", err)
		return "", err
	}

	// TODO:: Add handling where Firebase user exists and we should return
	// that existing user to be used for creation on our end, can be useful
	// in the event Firebase creation succeeds but fails on our end and needs
	// to be retried

	// TODO:: Migrate Firebase errors to service errors here

	return firebaseUser.UID, nil
}

func (s Service) GetFirebaseUserId(idToken string) (string, error) {
	ctx := context.Background()

	authClient, err := s.firebaseApp.Auth(ctx)
	if err != nil {
		log.Printf("encountered error while creating Firebase auth client due to %s", err)
		return "", err
	}

	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("unable to verify id token with Firebase due to %s", err)
		return "", err
	}

	return token.UID, nil
}

func (s Service) GetFirebaseUserEmail(firebaseUserId string) (email string, err error) {
	ctx := context.Background()

	authClient, err := s.firebaseApp.Auth(ctx)
	if err != nil {
		log.Printf("encountered error while creating Firebase auth client due to %s", err)
		return "", err
	}

	firebaseUser, err := authClient.GetUser(ctx, firebaseUserId)
	if err != nil {
		log.Printf("unable to fetch Firebase user due to %s", err)
		return "", err
	}

	return firebaseUser.Email, nil
}

// func (s Service) VerifyRecaptcha(recaptchaToken string) (err error) {
// 	if s.skipRecaptchaVerification {
// 		return nil
// 	}

// 	req, err := http.PostForm(s.recaptchaVerificationUrl, url.Values{
// 		"secret":   {s.recaptchaSecret},
// 		"response": {recaptchaToken},
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	defer req.Body.Close()

// 	respData, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		return
// 	}

// 	var recaptchaResponse GoogleRecaptchaResponse
// 	err = json.Unmarshal(respData, &recaptchaResponse)
// 	if err != nil {
// 		return
// 	}

// 	if !recaptchaResponse.Success {
// 		// https: //developers.google.com/recaptcha/docs/verify#error_code_reference
// 		for _, errorCode := range recaptchaResponse.ErrorCodes {
// 			if strings.Contains(errorCode, googleRecaptchaFailureKey) {
// 				return RecaptchaVerificationFailedError{}
// 			}
// 		}
// 		// TODO:: Unpack more response errors
// 		return errors.New(strings.Join(recaptchaResponse.ErrorCodes, "\n"))
// 	}

// 	return nil
// }

func (s Service) SendForgotPasswordEmail(email string) (err error) {
	// TODO:: Find implementation for backend reset as frontend
	// firebase library exposes email does not exist error

	// ctx := context.Background()

	// authClient, err := s.firebaseApp.Auth(ctx)
	// if err != nil {
	// 	log.Printf("encountered error while creating Firebase auth client due to %s", err)
	// 	return err
	// }

	// authClient.SendForgotPasswordEmail()

	return nil
}
