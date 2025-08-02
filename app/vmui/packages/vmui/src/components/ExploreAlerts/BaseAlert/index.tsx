import "./style.scss";
import { Alert as APIAlert } from "../../../types";
import { Link } from "react-router-dom";

interface BaseAlertProps {
  item: APIAlert;
}

const BaseAlert = ({ item }: BaseAlertProps) => {
  const query = item?.expression;

  const source = item.source;

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
              {(source && (
                <Link
                  to={source}
                  target="_blank"
                  rel="noreferrer"
                >
                  <pre>
                    <code className="language-promql">{query}</code>
                  </pre>
                </Link>
              )) || (
                <pre>
                  <code className="language-promql">{query}</code>
                </pre>
              )}
            </td>
          </tr>
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
    </div>
  );
};

export default BaseAlert;
