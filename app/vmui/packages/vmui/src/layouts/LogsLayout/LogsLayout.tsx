import { FC, useEffect } from "preact/compat";
import Header from "../Header/Header";
import { Outlet, useLocation } from "react-router-dom";
import "./style.scss";
import { getAppModeEnable } from "../../utils/app-mode";
import classNames from "classnames";
import Footer from "../Footer/Footer";
import router, { routerOptions } from "../../router";
import useDeviceDetect from "../../hooks/useDeviceDetect";
import ControlsLogsLayout from "./ControlsLogsLayout";
import { footerLinksToLogs } from "../../constants/footerLinks";
import useFetchFlags from "../../hooks/useFetchFlags";

const LogsLayout: FC = () => {
  const appModeEnable = getAppModeEnable();
  const { isMobile } = useDeviceDetect();
  const { pathname } = useLocation();

  const setDocumentTitle = () => {
    const defaultTitle = "UI for VictoriaLogs";
    const routeTitle = routerOptions[router.home]?.title;
    document.title = routeTitle ? `${routeTitle} - ${defaultTitle}` : defaultTitle;
  };

  useEffect(setDocumentTitle, [pathname]);
  useFetchFlags();

  return <section className="vm-container">
    <Header controlsComponent={ControlsLogsLayout}/>
    <div
      className={classNames({
        "vm-container-body": true,
        "vm-container-body_mobile": isMobile,
        "vm-container-body_app": appModeEnable
      })}
    >
      <Outlet/>
    </div>
    {!appModeEnable && <Footer links={footerLinksToLogs}/>}
  </section>;
};

export default LogsLayout;
