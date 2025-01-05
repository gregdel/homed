const headers = (method, data) => {
  let h = {
    method: method,
    headers: {
      "Content-Type": "application/json",
    },
  };

  if (data !== null) {
    if (typeof data === "string") {
      h.body = data;
    } else {
      h.body = JSON.stringify(data);
    }
  }

  return h;
};

export function request(
  eventPrefix,
  method,
  url,
  data = null,
  callbackEvents = null,
  mainPayload = null
) {
  // Events
  const pending = `${eventPrefix}_PENDING`;
  const fulfilled = `${eventPrefix}_FULFILLED`;
  const errored = `${eventPrefix}_ERROR`;

  return function (dispatch) {
    dispatch({
      type: pending,
      payload: {
        main: mainPayload,
      },
    });

    const conf = headers(method, data);

    return fetch(url, conf)
      .then((response) => response.json())
      .then((body) => {
        if (body.status === "error") {
          dispatch({
            type: errored,
            payload: {
              response: body,
              main: mainPayload,
            },
          });
          return;
        }
        dispatch({
          type: fulfilled,
          payload: {
            response: body.data,
            main: mainPayload,
          },
        });
        if (callbackEvents) {
          for (let event of callbackEvents) {
            if (typeof event === "function") {
              event = event();
            }
            dispatch(event);
          }
        }
      })
      .catch((response) => {
        console.warn(response);
      });
  };
}
