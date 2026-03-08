type Toast = {
  id: number;
  message: string;
};

function createToastStore() {
  let toasts = $state<Toast[]>([]);
  let nextId = 0;

  return {
    get toasts() {
      return toasts;
    },
    showError(message: string) {
      const id = nextId++;
      toasts.push({ id, message });
      setTimeout(() => {
        const index = toasts.findIndex((t) => t.id === id);
        if (index !== -1) toasts.splice(index, 1);
      }, 3000);
    },
    dismiss(id: number) {
      const index = toasts.findIndex((t) => t.id === id);
      if (index !== -1) toasts.splice(index, 1);
    },
  };
}

export const toastStore = createToastStore();