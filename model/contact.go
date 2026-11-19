package model

type Contact struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Mobile   string `json:"mobile"`
    Home     string `json:"home"`
}