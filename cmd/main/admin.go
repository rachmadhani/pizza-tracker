package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LoginData struct {
	Error string
}

type AdminDashboardData struct {
	Username string
}

type LoginForm struct {
	Username string `form:"username" binding:"required,min=3,max=50"`
	Password string `form:"password" binding:"required,min=6"`
}

func (h *Handler) HandleLoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", LoginData{})
}

func (h *Handler) HandleLoginPost(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{
			Error: "Invalid Input: " + err.Error(),
		})
		return
	}

	user, err := h.users.AuthenticateUser(form.Username, form.Password)
	if err != nil {
		c.HTML(http.StatusOK, "login.tmpl", LoginData{
			Error: "Invalid Credentials",
		})
		return
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	slog.Info("User logged in", "user_id", 42, "action", "login")

	SetSessionValue(c, "userID", fmt.Sprintf("%v", user.ID))
	SetSessionValue(c, "username", user.Username)

	c.Redirect(http.StatusSeeOther, "/admin")
}

func (h *Handler) HandleLogout(c *gin.Context) {
	if err := ClearSession(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func (h *Handler) ServeAdminDashboard(c *gin.Context) {
	username := GetSessionString(c, "username")

	c.HTML(http.StatusOK, "admin.tmpl", AdminDashboardData{
		Username: username,
	})
}
