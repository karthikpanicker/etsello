package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type apiHandler struct {
	handlerCom     *handlerCommon
	trelloManager  appDataManager
	gTManager      appDataManager
	todoistManager appDataManager
	dCache         dataStore
}

type trialInfo struct {
	SelectedBoardID string `json:"boardId"`
	SelectedListID  string `json:"listId"`
}

func newAPIHandler(cache dataStore) *apiHandler {
	ah := new(apiHandler)
	ah.handlerCom = newHandlerCommon()
	ah.trelloManager = getAppManager("trello")
	ah.gTManager = getAppManager("gtask")
	ah.todoistManager = getAppManager("todoist")
	ah.dCache = cache
	return ah
}

func (ah *apiHandler) getBordLists(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, _ := ah.dCache.getUserInfo(userID)
	params := mux.Vars(r)
	boardID := params["boardId"]
	if boardID == "" {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	requestParams := make(map[string]interface{})
	requestParams[trelloBoardIDKey] = boardID
	boardLists, _ := ah.trelloManager.getAppData(info, trelloBoardListRequest, requestParams)
	ah.handlerCom.ProcessResponse(boardLists, w)
}

func (ah *apiHandler) saveTrelloConfiguration(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, _ := ah.dCache.getUserInfo(userID)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&info.TrelloDetails)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	if info.TrelloDetails.SelectedBoardID == "" || info.TrelloDetails.SelectedListID == "" {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardList, w)
		return
	}
	// Mark the selected board in the list of boards as selected and mark others as unselected
	for index, board := range info.TrelloDetails.TrelloBoards {
		if board.ID == info.TrelloDetails.SelectedBoardID {
			info.TrelloDetails.TrelloBoards[index].IsSelected = true
		} else {
			info.TrelloDetails.TrelloBoards[index].IsSelected = false
		}
	}
	info.TrelloDetails.FromDate = ah.setFromDate(info.TrelloDetails.TransactionFilter)
	ah.dCache.saveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedTrello, w)
}

func (ah *apiHandler) getGTasksLists(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, _ := ah.dCache.getUserInfo(userID)
	tasks, err := ah.gTManager.getAppData(info, gTaskGetListsRequest, nil)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(err.Error(), w)
	}
	ah.handlerCom.ProcessResponse(tasks, w)
}

func (ah *apiHandler) saveGTasksConfig(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, _ := ah.dCache.getUserInfo(userID)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&info.GTasksDetails)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(messageInvalidBoardID, w)
		return
	}
	info.GTasksDetails.FromDate = ah.setFromDate(info.GTasksDetails.TransactionFilter)
	ah.dCache.saveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedGTasks, w)
}

func (ah *apiHandler) getTodoistProjects(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, _ := ah.dCache.getUserInfo(userID)
	projects, err := ah.todoistManager.getAppData(info, todoistProjectsRequest, nil)
	if err != nil {
		ah.handlerCom.ProcessErrorMessage(err.Error(), w)
	}
	ah.handlerCom.ProcessResponse(projects, w)
}
func (ah *apiHandler) saveTodoistConfig(w http.ResponseWriter, r *http.Request) {
	userID := ah.handlerCom.GetValueForKeyFromSession(r, userID).(int)
	info, _ := ah.dCache.getUserInfo(userID)
	decoder := json.NewDecoder(r.Body)
	_ = decoder.Decode(&info.TodoistDetails)
	info.TodoistDetails.FromDate = ah.setFromDate(info.TodoistDetails.TransactionFilter)
	ah.dCache.saveDetailsToCache(userID, *info)
	ah.handlerCom.ProcessSuccessMessage(messageSavedGTasks, w)
}

func (ah *apiHandler) setFromDate(filter int) int {
	fromDate := 0
	switch filter {
	case 1:
		fromDate = int(time.Now().Unix())
		break
	case 2:
		fromDate = int(time.Now().AddDate(0, -1, 0).Unix())
		break
	case 3:
		fromDate = int(time.Now().AddDate(0, 0, -7).Unix())
		break
	case 4:
		fromDate = 0
	default:
		fromDate = 0
	}
	return fromDate
}
