import { getEndpointPrefix } from "./config-helper.js";

export function decodeDate(s) {
  const dtFormat = new Intl.DateTimeFormat("pt-BR", {
    // timeStyle: "long",
    dateStyle: "long",
    timeZone: "America/Sao_Paulo",
  });

  return dtFormat.format(new Date(s * 1000));
}

export function decodeTime(s) {
  const dtFormat = new Intl.DateTimeFormat("pt-BR", {
    timeStyle: "long",
    // dateStyle: "full",
    timeZone: "America/Sao_Paulo",
  });

  return dtFormat.format(new Date(s * 1000));
}

export function getAvailableCenters() {
  return new Promise((resolve, reject) => {
    fetch(getEndpointPrefix() + "/api/v1/centers").then((response) =>
      response.json().then((centers) => resolve(centers))
    );
  });
}

export function getAvailableUsers(preProcessingFunc = null) {
  return new Promise((resolve, reject) => {
    fetch(getEndpointPrefix() + "/api/v1/users").then((response) =>
      response.json().then((users) => {
        if (preProcessingFunc) preProcessingFunc(users);
        resolve(users);
      })
    );
  });
}

// Function that calls the /api/v1/user_sequence endpoint to create a new UserSequence
export function createUserSequence(userSequence) {
  return new Promise((resolve, reject) => {
    fetch(getEndpointPrefix() + "/api/v1/user_sequence", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userSequence),
    }).then((response) =>
      response.json().then((userSequence) => resolve(userSequence))
    );
  });
}

// Function that calls the /api/v1/user_sequences endpoint to fetch UserSequences
// If processID is provided, it will use it as a URL query parameter to only fetch the UserSequences with such processID
export function getUserSequences(processID) {
  return new Promise((resolve, reject) => {
    let url = getEndpointPrefix() + "/api/v1/user_sequences";
    if (processID) url += `?processID=${processID}`;
    fetch(url).then((response) =>
      response.json().then((userSequences) => resolve(userSequences))
    );
  });
}

export function getNameFromUser(user) {
  return `${user.FirstName} ${user.LastName}`;
}

// Function that calls the /api/v1/user_sequence/count_completion POST endpoint,
// passing `processID`  as a query parameter
export function countCompletion(processID) {
  return new Promise((resolve, reject) => {
    fetch(
      `${getEndpointPrefix()}/api/v1/user_sequence/count_completion?processID=${processID}`,
      {
        method: "POST",
      }
    ).then((response) =>
      response.json().then((updatedSequence) => resolve(updatedSequence))
    );
  });
}

// Function that calls the /api/v1/process/:process_id/set_status/:status PUT endpoint to set the status of a `Process`
export function setProcessStatus(processID, status) {
  return new Promise((resolve, reject) => {
    fetch(
      `${getEndpointPrefix()}/api/v1/process/${processID}/set_status/${status}`,
      {
        method: "PUT",
      }
    ).then((response) =>
      response.json().then((updatedProcess) => resolve(updatedProcess))
    );
  });
}
// Function that calls the /api/v1/current_user/permissions GET endpoint to get the permissions of the current user
export function getCurrentUserPermissions() {
  return new Promise((resolve, reject) => {
    fetch(getEndpointPrefix() + "/api/v1/current_user/permissions").then(
      (response) =>
        response.json().then((permissions) => {
          resolve(permissions);
        })
    );
  });
}

export function hasPermission(permissions, targetID) {
  let _filteredVec = permissions.filter((perm) => perm.ID == targetID);
  console.log("[hasPermission]: {permissions, targetID, _filteredVec} = ", {
    permissions,
    targetID,
    _filteredVec,
  });
  return _filteredVec.length > 0;
}

export async function coreProcessUpdater(
  resolve,
  reject,
  set,
  modifiable = false
) {
  let rows = [];
  let promises = [];
  for (let queryParam of set) {
    promises.push(
      new Promise((resolve, reject) => {
        let modifiableString = modifiable ? "true" : "false";
        fetch(
          getEndpointPrefix() +
            "/api/v1/processes?statusString=" +
            queryParam +
            "&onlyModifiableByUser=" +
            modifiableString
        ).then((response) =>
          response.json().then(function (processes) {
            console.log(processes);
            let processObj = {};
            console.log("[updateProcesses::processes]:", processes);
            for (const p of processes) {
              processObj = {};
              processObj.id = p.ID;
              processObj.assunto = p.Title;
              processObj.centro = p.Center.Name;
              processObj.tipo = p.ProcessType.Name;
              processObj.estado = p.ProcessStatus.Name;

              //TODO: GET THE FOLLOWING FROM THE DB!
              // processObj.pend = "Encaminhamento final"
              processObj.autor = getNameFromUser(p.User);

              console.log(
                "[updateProcesses]: Process just built: ",
                processObj
              );
              rows.push(processObj);
            }
            resolve();
          })
        );
      })
    );
  }
  await Promise.all(promises);
  resolve(rows);
}

export async function coreUploadFile(file, processID) {
  // https://javascript.info/formdata
  let formData = new FormData();
  formData.set("documents", file, file.name);
  formData.set("processID", processID);

  const response = fetch(`${getEndpointPrefix()}/api/v1/files`, {
    method: "POST",
    body: formData,
  });

  return response;
}

// Function that calls the /api/v1/minute endpoint to create a new Minute
export async function postNewMinute(
  content,
  description,
  processID,
  centerID,
  attachmentID,
  isIncoming = false,
  inboundProtocol = "",
  errorCallback,
  dispatchMethod,
  callback
) {
  console.log("[postNewMinute]: Entering");
  try {
    // Here we use the Number function to prevent `processID` from being converted to a string
    // TODO: Let the description be provided by the user
    let requestBody = JSON.stringify({
      Content: content,
      Description: description,
      ProcessID: processID,
      UserFileID: attachmentID,
      IsIncoming: isIncoming,
      InboundProtocol: inboundProtocol,
      CenterID: centerID,
    });

    const response = await fetch(getEndpointPrefix() + "/api/v1/minute", {
      method: "post",

      body: requestBody,

      headers: {
        "Content-type": "application/json; charset=UTF-8",
      },
    });

    console.log("[postNewMinute::requestBody]: ", requestBody);

    response.text().then((text) => {
      console.log(text);
      console.log("[postNewMinute]: Completed!");
    });
  } catch (err) {
    console.error(`Error: ${err}`);
    if (errorCallback) {
      errorCallback();
    }
    return;
  }

  if (callback) {
    callback();
  }
  if (dispatchMethod) {
    dispatchMethod("minuteWasPosted");
  }
}
