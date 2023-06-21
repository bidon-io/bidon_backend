export default function ({ path, message }) {
  return useFormSubmit({
    path,
    message,
    method: "post",
    hook: async (id) => await navigateTo(`${path}/${id}`),
  });
}
