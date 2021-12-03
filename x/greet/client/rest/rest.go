package rest

import (
	"fmt"
	"net/http"
	// "strconv"
	// "strings"
	"github.com/kava-labs/kava/x/greet/types"
	"github.com/cosmos/cosmos-sdk/client"
	// sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)


func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/greetings", getGreetingsHandler(cliCtx)).Methods("GET")
	
	r.HandleFunc(fmt.Sprintf("/greetings/{%s}",types.QueryGetGreeting), getGreetingHandler(cliCtx)).Methods("GET")
}
//

func getGreetingsHandler(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryListGreetings), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
// query 
func getGreetingHandler(cliCtx client.Context) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		gid := vars[types.QueryGetGreeting]
		var q = types.QueryGetGreetRequest{Id: gid}.Id
		bz, err := cliCtx.LegacyAmino.MarshalJSON(q)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/greet/%s", types.QueryGetGreeting), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
		}
}


func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
