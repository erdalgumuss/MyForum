package handlers

import (
	"fmt"
	"net/http"

	"MyForum/config"
	"MyForum/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

