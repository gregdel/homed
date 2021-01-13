import axios from "axios";

// This functions returns an axios instance, the token is added to the
// configuration if found in the localStorage
export function configureAxios(headers = {}) {
  // Get the token from the localStorate
  const token = localStorage.getItem("token");
  if (token) {
    headers = { Authorization: `Bearer ${token}` };
  }

  return axios.create({
    headers,
  });
}

// This function takes en event prefix to dispatch evens during the life of the
// request, it also take a promise (axios request)
export function request(
  eventPrefix,
  promise,
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

    return promise
      .then((response) => {
        if (response.data.status === "error") {
          dispatch({
            type: "ADD_ALERT_ERROR",
            payload: {
              message: response.data.message,
              main: mainPayload,
            },
          });
          dispatch({
            type: errored,
            payload: {
              response: response.data,
              main: mainPayload,
            },
          });
          return;
        }
        dispatch({
          type: fulfilled,
          payload: {
            response: response.data,
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
      .catch((error) => {
        // Unauthorized
        if (error.response && error.response.status == 401) {
          dispatch({
            type: "USER_LOGOUT",
          });
        }
        dispatch({
          type: "ADD_ALERT_ERROR",
          payload: {
            message: error.response.data,
            main: mainPayload,
          },
        });
      });
  };
}
