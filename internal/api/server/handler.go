package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlexeyBazhin/wbL0/internal/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	searchFormTmpl = []byte(`
	<html>
		<body>
		<form method="post" >
			<p><input type="text" name="uid"></p>
			<p><input type="submit" value="Получить заказ"></p>
		</form>
		</body>
	</html>
	`)
	searchErrFormTmpl = []byte(`
	<html>
		<body>
		<p>Введите корректный uid заказа!</p>
		%s
		</body>
	</html>
	`)
)

func (server *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write(searchFormTmpl)
		return
	}
	orderUid := r.FormValue("uid")
	if _, err := uuid.Parse(orderUid); err != nil {
		fmt.Fprintf(w, string(searchErrFormTmpl), searchFormTmpl)
	} else {
		http.Redirect(w, r, "/orders/"+orderUid, http.StatusFound)
	}
}

func (server *Server) orderHandler(w http.ResponseWriter, r *http.Request) {
	server.logger.Infof("[http-server] request: %v", r.URL.Path)
	vars := mux.Vars(r)
	orderUid, err := uuid.Parse(vars["order-uid"])
	if err != nil {
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
	}

	if modelJSON, err := server.redisClient.Get(server.ctx, orderUid.String()).Bytes(); err != nil {
		server.logger.Errorf("[http-server] get from Redis error: %w", err)
	} else {
		w.Write(modelJSON)
		server.logger.Info("[http-server] successful get from Redis")
		return
	}
	completeOrder, err := server.svc.GetOrderById(server.ctx, orderUid)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
		return
	}
	modelJSON := api.MakeJSONModel(completeOrder)
	if err := json.NewEncoder(w).Encode(&modelJSON); err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
		return
	}
	server.logger.Info("[http-server] successful get from PG")
}

// func (server *Server) StanListener(msg *stan.Msg) {
// 	model := &model.Model{}
// 	if err := json.Unmarshal(msg.Data, model); err != nil {
// 		server.logger.Error("[http-server] cannot unmarshal json model: ", zap.Error(err))
// 	}

// 	if order, err := server.createOrderService(model.Order); err != nil {
// 		server.logger.Error("[http-server] cannot create order: ", zap.Error(err))
// 	}

// }
