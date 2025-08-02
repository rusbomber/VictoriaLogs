import "./style.scss";
import { Rule as APIRule } from "../../../types";
import { Link, createSearchParams } from "react-router-dom";
import dayjs from "dayjs";
import { useNavigate } from "react-router-dom";

interface BaseRuleProps {
  item: APIRule;
}

const BaseRule = ({ item }: BaseRuleProps) => {
  const query = item?.query;
  const navigate = useNavigate();
  const openItemLink = (id: string) => {
    return () => {
      navigate({
        pathname: "/groups",
        search:   `group_id=${item.group_id}&alert_id=${id}`,
      });
    };
  };

  const openQueryLink = (time: string) => {
    return () => {
      const params = {
        "query": query,
        "g0.end_time": time,
      };
      navigate({
        pathname: "/",
        search: `?${createSearchParams(params).toString()}`,
      });
    };
  };

  const source = {
    pathname: "/",
    search: createSearchParams({
      "guery": query,
    }).toString(),
  };

  return (
    <div className="vm-explore-alerts-rule-item">
      <table>
        <tbody>
          <tr>
            <td className="vm-col-small">Name</td>
            <td>{item.name}</td>
          </tr>
          <tr>
            <td className="vm-col-small">Group</td>
            <td>
              <Link to={`/groups#group-${item.group_id}`}>{item.group_id}</Link>
            </td>
          </tr>
          <tr>
            <td className="vm-col-small">Query</td>
            <td>
              <Link to={source}>
                <pre>
                  <code className="language-promql">{query}</code>
                </pre>
              </Link>
            </td>
          </tr>
          <tr>
            <td className="vm-col-small">Duration</td>
            <td>{dayjs.duration(item.duration, "seconds").humanize()}</td>
          </tr>
          <tr>
            <td className="vm-col-small">Last evaluation</td>
            <td>{dayjs().to(dayjs(item.lastEvaluation))}</td>
          </tr>
          {item.lastError && (
            <tr>
              <td className="vm-col-small">Last error</td>
              <td>{item.lastError}</td>
            </tr>
          )}
          {item.labels && (
            <tr>
              <td className="vm-col-small">Labels</td>
              <td>
                <div className="vm-badge-container">
                  {Object.entries(item.labels).map(([name, value]) => (
                    <span
                      key={name}
                      className="vm-badge"
                    >{`${name}: ${value}`}</span>
                  ))}
                </div>
              </td>
            </tr>
          )}
        </tbody>
      </table>
      {item.annotations && (
        <>
          <span className="title">Annotations</span>
          <table>
            <tbody>
              {Object.entries(item.annotations || {}).map(([name, value]) => (
                <tr key={name}>
                  <td className="vm-col-small">{name}</td>
                  <td>{value}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
      {item.updates && (
        <>
          <span className="title">{`Last updates ${item.updates.length}/20`}</span>
          <table>
            <thead>
              <tr>
                <th className="vm-col-small">Updated at</th>
                <th className="vm-col-small">Series returned</th>
                <th className="vm-col-small">Series fetched</th>
                <th className="vm-col-small">Duration</th>
                <th className="vm-col-small">Executed at</th>
              </tr>
            </thead>
            <tbody>
              {item.updates.map((update) => (
                <tr
                  key={update.at}
                  className="hoverable"
                  onClick={openQueryLink(update.at)}
                >
                  <td className="vm-col-small">{update.time}</td>
                  <td className="vm-col-small">{update.samples}</td>
                  <td className="vm-col-small">{update.series_fetched}</td>
                  <td className="vm-col-small">{dayjs.duration(update.duration / 1e6).asSeconds().toFixed(3)}s</td>
                  <td className="vm-col-small">{update.at}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
      {item.alerts && (
        <>
          <span className="title">Alerts</span>
          <table>
            <thead>
              <tr>
                <th className="vm-col-small">Active since</th>
                <th className="vm-col-small">Value</th>
                <th>Labels</th>
              </tr>
            </thead>
            <tbody>
              {item.alerts.map((alert) => (
                <tr
                  id={`alert-${alert.id}`}
                  key={alert.id}
                  className="hoverable"
                  onClick={openItemLink(alert.id)}
                >
                  <td className="vm-col-small">
                    {dayjs().to(dayjs(alert.activeAt))}
                  </td>
                  <td className="vm-col-small">
                    <div className="vm-badge">{alert.value}</div>
                  </td>
                  <td className="vm-badge-container">
                    {Object.entries(alert.labels || {}).map(([name, value]) => (
                      <span
                        key={name}
                        className="vm-badge"
                      >{`${name}: ${value}`}</span>
                    ))}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  );
};

export default BaseRule;
