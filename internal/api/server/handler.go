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
	server.logger.Infof("[http-server] request: %v", r.URL.Path)

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
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderUid, err := uuid.Parse(vars["order-uid"])
	if err != nil {
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
	}

	if byteFromCache, err := server.svc.PullFromCache(server.ctx, orderUid); err != nil {
		server.logger.Errorf("[http-server] get from Redis error: %w", err)
	} else {
		// var modelJSON api.ModelJSON
		// if err := json.Unmarshal(byteFromCache, &modelJSON); err != nil {
		// 	fmt.Println(fmt.Errorf("failed to unmarshal json data: %s", err.Error()))
		// }
		// modelByte, _ := json.MarshalIndent(modelJSON, "", " ")
		// w.Write(modelByte)
		w.Write(byteFromCache) //раньше во время пуша в рэдис внутри stanListener не применялся MarshalIndent - поэтому нужен был закомментированный выше код

		server.logger.Info("[http-server] successful get from Redis")
		return
	}

	//если не удалось получить из кэша
	completeOrder, err := server.svc.GetOrderById(server.ctx, orderUid)
	if err != nil {
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
		return
	}
	modelJSON := api.MakeJSONModel(completeOrder)
	if err != nil {
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
		return
	}
	server.logger.Info("[http-server] successful get from PG")

	if modelByte, err := json.MarshalIndent(modelJSON, "", " "); err != nil {
		json.NewEncoder(w).Encode(&api.ErrorJSON{
			Err: err.Error(),
		})
		return
	} else {
		w.Write(modelByte)
	}

	//полученные данные нужно записать в кэш
	modelByte, err := json.MarshalIndent(modelJSON, "", " ")
	if err != nil {
		server.logger.Errorf("[http-server] failed to marshal modelJSON: %w", err)
		return
	}
	if err := server.svc.PushToCache(
		server.ctx, completeOrder.Order.Id, modelByte,
	); err != nil {
		server.logger.Errorf("[http-server] failed to cache data: %w", err)
		return
	}
	server.logger.Info("[stan-listener] DATA SAVED TO CACHE")
}
