import { FC, useMemo, useState } from "preact/compat";
import "./style.scss";
import "uplot/dist/uPlot.min.css";
import { AlignedData } from "uplot";
import { TimeParams } from "../../../types";
import classNames from "classnames";
import { LogHits } from "../../../api/types";
import { GraphOptions, GRAPH_STYLES } from "./types";
import BarHitsOptions from "./BarHitsOptions/BarHitsOptions";
import BarHitsPlot from "./BarHitsPlot/BarHitsPlot";
import { calculateTotalHits } from "../../../utils/logs";
import { LOGS_GROUP_BY, LOGS_LIMIT_HITS } from "../../../constants/logs";
import { getDurationFromMilliseconds } from "../../../utils/time";

interface Props {
  logHits: LogHits[];
  data: AlignedData;
  period: TimeParams;
  durationMs?: number;
  setPeriod: ({ from, to }: { from: Date, to: Date }) => void;
  onApplyFilter: (value: string) => void;
}

const BarHitsChart: FC<Props> = ({ logHits, data: _data, period, setPeriod, onApplyFilter, durationMs }) => {
  const [graphOptions, setGraphOptions] = useState<GraphOptions>({
    graphStyle: GRAPH_STYLES.BAR,
    stacked: false,
    fill: false,
    hideChart: false,
  });

  const totalHits = useMemo(() => calculateTotalHits(logHits), [logHits]);

  return (
    <div
      className={classNames({
        "vm-bar-hits-chart__wrapper": true,
        "vm-bar-hits-chart__wrapper_hidden": graphOptions.hideChart
      })}
    >
      <div className="vm-bar-hits-chart-header">
        <div className="vm-bar-hits-chart-header-info">
          <p>Top <b>{LOGS_LIMIT_HITS}</b> hits grouped by <b>{LOGS_GROUP_BY}</b></p>
          |
          <p>Total: <b>{totalHits.toLocaleString("en-US")}</b> hits</p>
          {durationMs !== undefined && (
            <>
              |
              <p>Duration: <b>{getDurationFromMilliseconds(durationMs)}</b></p>
            </>
          )}
        </div>
        <BarHitsOptions onChange={setGraphOptions}/>
      </div>
      {!graphOptions.hideChart && (
        <BarHitsPlot
          logHits={logHits}
          totalHits={totalHits}
          data={_data}
          period={period}
          setPeriod={setPeriod}
          onApplyFilter={onApplyFilter}
          graphOptions={graphOptions}
        />
      )}
    </div>
  );
};

export default BarHitsChart;
