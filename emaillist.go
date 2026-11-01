package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gocql/gocql"
)

func putNotificationSelect(w http.ResponseWriter, r *http.Request) { // Select all user mailing lists and active emails in them
	start := time.Now()

	type ResponseEMailListDataT struct {
		ID          string `json:"id"`
		EMail       string `json:"email"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	type ResponseEMailListT struct {
		ID          string                   `json:"id"`
		Name        string                   `json:"name"`
		Description string                   `json:"description"`
		Status      int                      `json:"status"`
		Owner       string                   `json:"owner"`
		Data        []ResponseEMailListDataT `json:"data"`
	}

	type ResponseEMailListAT struct {
		Data        []ResponseEMailListT `json:"data"`
		CurrentPage string               `json:"current"`
		NextPage    string               `json:"next"`
	}

	type SelectParam struct {
		CurrentPage string `json:"current"`
		PageSize    int    `json:"pagesize"`
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested a list of emaillist", claims["username"])

	if claims["role"] == "superadmin" || claims["role"] == "admin" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			ctx01 := context.Background()
			var scanner01 gocql.Scanner
			var iter *gocql.Iter

			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var selectParam SelectParam
			if err := json.Unmarshal(b, &selectParam); err != nil {
				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			pageSize := max(selectParam.PageSize, 0)

			var responseEmaiListData ResponseEMailListDataT
			var responseEMailList ResponseEMailListT
			var responseEMailListA ResponseEMailListAT
			var pageState []byte

			if selectParam.CurrentPage != "" {
				var err error
				pageState, err = base64.StdEncoding.DecodeString(selectParam.CurrentPage)
				if err != nil {
					log.Print("CurrentPage decode ERROR (" + err.Error() + ")!")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			switch claims["role"] {
			case "superadmin":
				if selectParam.PageSize == 0 {
					scanner01 = Session.Query(`SELECT id, name, description, status, user_id FROM emaillist`).WithContext(ctx01).Consistency(ConsistencyRead).Iter().Scanner()
				} else {
					iter = Session.Query(`SELECT id, name, description, status, user_id FROM emaillist`).PageSize(pageSize).PageState(pageState).Iter()
				}
			case "admin":
				if selectParam.PageSize == 0 {
					scanner01 = Session.Query(`SELECT id, name, description, status, user_id FROM emaillist WHERE user_id = ?`, claims["userid"]).WithContext(ctx01).Consistency(ConsistencyRead).Iter().Scanner()
				} else {
					iter = Session.Query(`SELECT id, name, description, status, user_id FROM emaillist WHERE user_id = ?`, claims["userid"]).PageSize(pageSize).PageState(pageState).Iter()
				}
			}

			if selectParam.PageSize > 0 {
				nextPageState := iter.PageState()
				if len(nextPageState) > 0 {
					responseEMailListA.NextPage = base64.StdEncoding.EncodeToString(nextPageState)
				} else {
					responseEMailListA.NextPage = ""
				}
				responseEMailListA.CurrentPage = selectParam.CurrentPage

				scanner01 = iter.Scanner()
			}

			for scanner01.Next() {
				var emailListID string
				var emailListName string
				var emailListDescription string
				var emailListStatus int
				var emailListUserID gocql.UUID
				err := scanner01.Scan(&emailListID, &emailListName, &emailListDescription, &emailListStatus, &emailListUserID)
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/notification/select", "500", "GET", "getEMailListSelect").Set(duration)

					log.Print(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {

					var emailListUserName string
					ctx012 := context.Background()
					scanner012 := Session.Query(`SELECT user_name FROM user_id WHERE user_id = ?`, emailListUserID).WithContext(ctx012).Consistency(ConsistencyRead).Iter().Scanner()
					for scanner012.Next() {
						err := scanner012.Scan(&emailListUserName)
						if err != nil {
							log.Print(err)
						}
					}

					var responseEMailListDataA []ResponseEMailListDataT
					ctx02 := context.Background()
					scanner02 := Session.Query(`SELECT id, email, description, status FROM emaillist_data WHERE emaillist_id = ?`, emailListID).WithContext(ctx02).Consistency(ConsistencyRead).Iter().Scanner()
					for scanner02.Next() {
						var emailListDataID string
						var email string
						var description string
						var status int
						err := scanner02.Scan(&emailListDataID, &email, &description, &status)
						if err != nil {
							duration := time.Since(start).Seconds()
							httpDuration.WithLabelValues("/notification/select", "500", "GET", "getEMailListSelect").Set(duration)

							log.Print(err)
							w.WriteHeader(http.StatusInternalServerError)
							return
						} else {
							responseEmaiListData.ID = emailListDataID
							responseEmaiListData.EMail = email
							responseEmaiListData.Description = description
							responseEmaiListData.Status = status

							responseEMailListDataA = append(responseEMailListDataA, responseEmaiListData)
						}
					}
					responseEMailList.ID = emailListID
					responseEMailList.Name = emailListName
					responseEMailList.Description = emailListDescription
					responseEMailList.Status = emailListStatus
					responseEMailList.Data = responseEMailListDataA
					responseEMailList.Owner = emailListUserName

					responseEMailListA.Data = append(responseEMailListA.Data, responseEMailList)
				}
			}

			responseJSON, err := json.Marshal(responseEMailListA)
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/select", "500", "GET", "getEMailListSelect").Set(duration)

				log.Print("JSON MARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/select", "200", "GET", "getEMailListSelect").Set(duration)

				w.Header().Set("Content-Type", "application/json")
				w.Write(responseJSON)
				return
			}
		} else {
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/notification/select", "500", "GET", "getEMailListSelect").Set(duration)

			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/notification/select", "403", "GET", "getEMailListSelect").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func putNotificationCreate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested the creation of a emaillis!", claims["username"])

	type EMailList struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if claims["role"] == "superadmin" || claims["role"] == "admin" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/create", "500", "PUT", "putNotificationCreate").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			id, err := gocql.RandomUUID()
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/create", "500", "PUT", "putNotificationCreate").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var emailList EMailList
			if err := json.Unmarshal(b, &emailList); err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/create", "500", "PUT", "putNotificationCreate").Set(duration)

				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			{
				ctx := context.Background()
				err := Session.Query(`INSERT INTO emaillist (user_id, id, name, description, status) VALUES (?, ?, ?, ?, ?)`,
					claims["userid"], id, emailList.Name, emailList.Description, emailList.Status).WithContext(ctx).Exec()
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/notification/create", "500", "PUT", "putNotificationCreate").Set(duration)

					log.Print("Failed to create a emaillist (ERROR: ", err.Error(), ")!")
					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/notification/create", "200", "PUT", "putNotificationCreate").Set(duration)

					log.Printf("User %v successfully created emaillist!", claims["username"])
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		} else {
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/notification/create", "500", "PUT", "putNotificationCreate").Set(duration)

			log.Print("It was not possible to create a emaillist!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/notification/create", "403", "PUT", "putNotificationCreate").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func putNotificationEdit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested emaillist update!", claims["username"])

	type EMailListData struct {
		ID          string `json:"id"`
		EMail       string `json:"email"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	type EMailList struct {
		ID            string          `json:"id"`
		Name          string          `json:"name"`
		Description   string          `json:"description"`
		Status        int             `json:"status"`
		EMailListData []EMailListData `json:"data"`
	}

	if claims["role"] == "superadmin" || claims["role"] == "admin" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/edit", "500", "PUT", "putNotificationEdit").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx := context.Background()

			var emailList EMailList
			if err := json.Unmarshal(b, &emailList); err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/edit", "500", "PUT", "putNotificationEdit").Set(duration)

				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = Session.Query(`DELETE FROM emaillist_data WHERE emaillist_id = ?`, emailList.ID).WithContext(ctx).Exec()
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/edit", "500", "PUT", "putNotificationEdit").Set(duration)

				log.Print("Failed to remove emaillist (ERROR: ", err.Error(), ")!")
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				err = Session.Query(`UPDATE emaillist SET name = ?, description = ?, status = ? WHERE user_id = ? AND id = ?`,
					emailList.Name, emailList.Description, emailList.Status, claims["userid"], emailList.ID).WithContext(ctx).Exec()
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/notification/edit", "500", "PUT", "putNotificationEdit").Set(duration)

					log.Print("Failed to update emaillist (ERROR: ", err.Error(), ")!")
					w.WriteHeader(http.StatusInternalServerError)
					return
				} else {
					for _, v := range emailList.EMailListData {
						uuidEMailListData, err := gocql.RandomUUID()
						if err != nil {
							duration := time.Since(start).Seconds()
							httpDuration.WithLabelValues("/notification/edit", "500", "PUT", "putNotificationEdit").Set(duration)

							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						err = Session.Query(`INSERT INTO emaillist_data (emaillist_id, id, email, description, status) VALUES (?, ?, ?, ?, ?)`,
							emailList.ID, uuidEMailListData, v.EMail, v.Description, v.Status).WithContext(ctx).Exec()
						if err != nil {
							duration := time.Since(start).Seconds()
							httpDuration.WithLabelValues("/notification/edit", "500", "PUT", "putNotificationEdit").Set(duration)

							log.Print("Failed to update emaillist data (ERROR: ", err.Error(), ")!")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					}
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/notification/edit", "200", "PUT", "putNotificationEdit").Set(duration)

					log.Printf("User %v successfully created emaillist!", claims["username"])
					w.WriteHeader(http.StatusOK)
					return
				}
			}
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/notification/edit", "403", "PUT", "putNotificationEdit").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}

func putNotificationRemove(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Printf("User %v requested emaillist removal!", claims["username"])

	type EMailList struct {
		ID string `json:"id"`
	}

	if claims["role"] == "superadmin" || claims["role"] == "admin" {
		versionAPI := chi.URLParam(r, "version")
		if versionAPI == "1" {
			b, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/remove", "500", "PUT", "putNotificationRemove").Set(duration)

				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var emailList EMailList
			if err := json.Unmarshal(b, &emailList); err != nil {
				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/remove", "500", "PUT", "putNotificationRemove").Set(duration)

				log.Print("JSON UNMARSHAL ERROR (" + err.Error() + ")!")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			{
				ctx := context.Background()
				err = Session.Query(`
				BEGIN BATCH
					DELETE FROM emaillist_data WHERE emaillist_id = ?
					DELETE FROM emaillist WHERE user_id = ? AND id = ?
				APPLY BATCH`,
					emailList.ID, claims["userid"], emailList.ID).WithContext(ctx).Exec()
				if err != nil {
					duration := time.Since(start).Seconds()
					httpDuration.WithLabelValues("/notification/remove", "500", "PUT", "putNotificationRemove").Set(duration)

					log.Print("Failed to remove emaillist (ERROR: ", err.Error(), ")!")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				duration := time.Since(start).Seconds()
				httpDuration.WithLabelValues("/notification/remove", "200", "PUT", "putNotificationRemove").Set(duration)

				log.Printf("User %v successfully deleted emaillist!", claims["username"])
				w.WriteHeader(http.StatusOK)
				return
			}
		} else {
			duration := time.Since(start).Seconds()
			httpDuration.WithLabelValues("/notification/remove", "500", "PUT", "putNotificationRemove").Set(duration)

			log.Print("Failed to remove emaillist (API VERSION ERROR)!")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		duration := time.Since(start).Seconds()
		httpDuration.WithLabelValues("/notification/remove", "403", "PUT", "putNotificationRemove").Set(duration)

		w.WriteHeader(http.StatusForbidden)
		return
	}
}
