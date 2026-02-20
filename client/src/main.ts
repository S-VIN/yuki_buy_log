import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

// ── iOS Safari: предотвращаем прокрутку страницы при открытии клавиатуры ─────
// Сбрасываем позицию если window всё же прокрутился
window.addEventListener('scroll', () => {
  if (window.scrollX !== 0 || window.scrollY !== 0) {
    window.scrollTo(0, 0);
  }
}, { passive: true });

// visualViewport scroll (Safari >= 13 при появлении клавиатуры)
if (window.visualViewport) {
  window.visualViewport.addEventListener('scroll', () => {
    window.scrollTo(0, 0);
  });
}

// Запрещаем touchmove на элементах которые не имеют собственного скролла
document.addEventListener('touchmove', (e) => {
  let el = e.target as HTMLElement | null;
  while (el && el !== document.body) {
    const style = window.getComputedStyle(el);
    const overflow = style.overflow + style.overflowY;
    if (/(auto|scroll)/.test(overflow) && el.scrollHeight > el.clientHeight) {
      return; // разрешаем скролл внутри прокручиваемых контейнеров
    }
    el = el.parentElement;
  }
  e.preventDefault();
}, { passive: false });
// ─────────────────────────────────────────────────────────────────────────────

const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
