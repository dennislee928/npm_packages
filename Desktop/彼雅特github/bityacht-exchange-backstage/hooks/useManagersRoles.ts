export default function useManagersRoles() {
  const roles = [
    {
      id: 1,
      name: {
        zh: "超級帳號",
        en: "Administrator",
      },
    },
    {
      id: 2,
      name: {
        zh: "法遵",
        en: "Compliance",
      },
    },
    {
      id: 3,
      name: {
        zh: "客服",
        en: "Customer Service",
      },
    },
    {
      id: 4,
      name: {
        zh: "風控",
        en: "Risk Management",
      },
    },
    {
      id: 5,
      name: {
        zh: "財務",
        en: "Finance",
      },
    },
  ];

  return {
    getRoleNameFromId: (id: number) =>
      roles.find((role) => role.id === id)?.name,
    getRoleIdFromName: (name: string) =>
      roles.find((role) => role.name.en === name || role.name.zh === name)?.id,
  };
}
