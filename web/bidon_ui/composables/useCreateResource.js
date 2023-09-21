export default function ({
  path,
  message,
  onError = async () => {
    /* no operation function */
  },
}) {
  return useFormSubmit({
    path,
    message,
    method: "post",
    onSuccess: async (id) => await navigateTo(`${path}/${id}`),
    onError: onError,
  });
}
