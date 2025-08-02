export const getGroupsUrl = (server: string, ruleType: string): string => {
  let groupUrl = `${server}/vmalert/api/v1/rules?datasource_type=vlogs`;
  if (ruleType) {
    groupUrl = `${groupUrl}&type=${ruleType}`;
  }
  return groupUrl;
};

export const getItemUrl = (
  server: string,
  groupId: string,
  id: string,
  mode: string,
): string => {
  return `${server}/vmalert/api/v1/${mode}?group_id=${groupId}&${mode}_id=${id}`;
};

export const getNotifiersUrl = (server: string): string => {
  return `${server}/vmalert/api/v1/notifiers`;
};
