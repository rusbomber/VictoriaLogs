import { AppStateProvider } from "../state/common/StateContext";
import { TimeStateProvider } from "../state/time/TimeStateContext";
import { QueryStateProvider } from "../state/query/QueryStateContext";
import { LogsStateProvider } from "../state/logsPanel/LogsStateContext";
import { OverviewStateProvider } from "../state/overview/OverviewStateContext";
import { SnackbarProvider } from "./Snackbar";

import { combineComponents } from "../utils/combine-components";

const providers = [
  AppStateProvider,
  TimeStateProvider,
  QueryStateProvider,
  SnackbarProvider,
  LogsStateProvider,
  OverviewStateProvider,
];

export default combineComponents(...providers);
