import { FC, useMemo, useRef, useState, MouseEvent } from "preact/compat";
import classNames from "classnames";
import { Series } from "uplot";
import { LegendLogHits } from "../../../../api/types";
import { getStreamPairs } from "../../../../utils/logs";
import { formatNumberShort } from "../../../../utils/math";
import Popper from "../../../Main/Popper/Popper";
import useBoolean from "../../../../hooks/useBoolean";
import LegendHitsMenu from "../LegendHitsMenu/LegendHitsMenu";
import { ExtraFilter } from "../../../../pages/OverviewPage/FiltersBar/types";

interface Props {
  legend: LegendLogHits;
  series: Series[];
  onRedrawGraph: () => void;
  onApplyFilter: (value: ExtraFilter) => void;
}

const BarHitsLegendItem: FC<Props> = ({ legend, series, onRedrawGraph, onApplyFilter }) => {
  const {
    value: openContextMenu,
    setTrue: handleOpenContextMenu,
    setFalse: handleCloseContextMenu,
  } = useBoolean(false);

  const legendRef = useRef<HTMLDivElement>(null);
  const [clickPosition, setClickPosition] = useState<{ top: number; left: number } | null>(null);

  const targetSeries = useMemo(() => series.find(s => s.label === legend.label), [series]);

  const fields = useMemo(() => getStreamPairs(legend.label), [legend.label]);

  const label = fields.join(", ");
  const totalShortFormatted = formatNumberShort(legend.total);

  const handleContextMenu = (e: MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    setClickPosition({ top: e.clientY, left: e.clientX });
    handleOpenContextMenu();
  };

  return (
    <div
      ref={legendRef}
      className={classNames({
        "vm-bar-hits-legend-item": true,
        "vm-bar-hits-legend-item_other": legend.isOther,
        "vm-bar-hits-legend-item_active": openContextMenu,
        "vm-bar-hits-legend-item_hide": !targetSeries?.show,
      })}
      onClick={handleContextMenu}
    >
      <div
        className="vm-bar-hits-legend-item__marker"
        style={{ backgroundColor: `${legend.stroke}` }}
      />
      <div className="vm-bar-hits-legend-item__label">{label}</div>
      <span className="vm-bar-hits-legend-item__total">({totalShortFormatted})</span>

      <Popper
        placement="fixed"
        open={openContextMenu}
        buttonRef={legendRef}
        placementPosition={clickPosition}
        onClose={handleCloseContextMenu}
      >
        <LegendHitsMenu
          legend={legend}
          fields={fields}
          series={series}
          onApplyFilter={onApplyFilter}
          onRedrawGraph={onRedrawGraph}
          onClose={handleCloseContextMenu}
        />
      </Popper>
    </div>
  );
};

export default BarHitsLegendItem;
