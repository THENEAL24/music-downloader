const API = {
  streamTxt:    '/api/stream/txt',
  streamYandex: '/api/stream/yandex',
  result:       '/api/result/',
};

// ── DOM-ссылки ────────────────────────────────────────────────────────────

const $ = id => document.getElementById(id);

const dom = {
  form:            $('main-form'),
  tabTxt:          $('tab-txt'),
  tabYandex:       $('tab-yandex'),
  sectionTxt:      $('section-txt'),
  sectionYandex:   $('section-yandex'),
  dropZone:        $('drop-zone'),
  fileInput:       $('file-input'),
  fileName:        $('file-name'),
  yandexUrl:       $('yandex-url'),
  yandexToken:     $('yandex-token'),
  tracklistWrap:   $('tracklist-preview'),
  tracklistCount:  $('tracklist-count'),
  tracklistUl:     $('tracklist-ul'),
  optWorkers:      $('opt-workers'),
  optLimit:        $('opt-limit'),
  submitBtn:       $('submit-btn'),
  progressSection: $('progress-section'),
  progressLabel:   $('progress-label'),
  progressPct:     $('progress-pct'),
  progressBar:     $('progress-bar'),
  progressWrap:    document.querySelector('.progress-bar-wrap'),
  log:             $('log'),
  stats:           $('stats'),
  statOk:          $('stat-ok'),
  statFail:        $('stat-fail'),
  statPct:         $('stat-pct'),
  downloadBtn:     $('download-btn'),
  resetBtn:        $('reset-btn'),
};

// ── Состояние ─────────────────────────────────────────────────────────────

let activeTab = 'txt';

// ── Табы ──────────────────────────────────────────────────────────────────

function switchTab(tab) {
  activeTab = tab;
  dom.tabTxt.classList.toggle('active', tab === 'txt');
  dom.tabYandex.classList.toggle('active', tab === 'yandex');
  dom.tabTxt.setAttribute('aria-selected', tab === 'txt');
  dom.tabYandex.setAttribute('aria-selected', tab === 'yandex');
  dom.sectionTxt.style.display    = tab === 'txt'    ? 'block' : 'none';
  dom.sectionYandex.style.display = tab === 'yandex' ? 'block' : 'none';
  dom.tracklistWrap.hidden = true;
}

// ── Drag & Drop ───────────────────────────────────────────────────────────

dom.dropZone.addEventListener('dragover', e => {
  e.preventDefault();
  dom.dropZone.classList.add('drag-over');
});
dom.dropZone.addEventListener('dragleave', () => dom.dropZone.classList.remove('drag-over'));
dom.dropZone.addEventListener('drop', e => {
  e.preventDefault();
  dom.dropZone.classList.remove('drag-over');
  const file = e.dataTransfer.files[0];
  if (file) setFile(file);
});
dom.fileInput.addEventListener('change', () => {
  if (dom.fileInput.files[0]) setFile(dom.fileInput.files[0]);
});

function setFile(file) {
  dom.fileName.textContent = '📎 ' + file.name;

  const reader = new FileReader();
  reader.onload = e => {
    const lines = e.target.result
      .split(/\r?\n/)
      .map(l => l.trim())
      .filter(Boolean);
    renderTracklistPreview(lines);
  };
  reader.readAsText(file, 'UTF-8');
}

// ── Tracklist preview ─────────────────────────────────────────────────────

function renderTracklistPreview(lines) {
  const MAX_PREVIEW = 30;
  dom.tracklistCount.textContent = lines.length;
  dom.tracklistUl.innerHTML = '';

  const slice = lines.slice(0, MAX_PREVIEW);
  slice.forEach((t, i) => {
    const li = document.createElement('li');
    li.innerHTML = `<span>${i + 1}.</span> ${escHtml(t)}`;
    dom.tracklistUl.appendChild(li);
  });

  if (lines.length > MAX_PREVIEW) {
    const li = document.createElement('li');
    li.innerHTML = `<span>…ещё ${lines.length - MAX_PREVIEW} треков</span>`;
    dom.tracklistUl.appendChild(li);
  }

  dom.tracklistWrap.hidden = false;
}

// ── Log ───────────────────────────────────────────────────────────────────

function logAdd(cls, text) {
  const div = document.createElement('div');
  div.className = cls;
  div.textContent = text;
  dom.log.appendChild(div);
  dom.log.scrollTop = dom.log.scrollHeight;
}

function setProgress(done, total) {
  const pct = total > 0 ? Math.round(done / total * 100) : 0;
  dom.progressBar.style.width = pct + '%';
  dom.progressPct.textContent = pct + '%';
  dom.progressWrap.setAttribute('aria-valuenow', pct);
}

// ── SSE обработка сообщений ───────────────────────────────────────────────

function handleMsg(msg, total) {
  switch (msg.type) {
    case 'start':
      dom.progressLabel.textContent = `Скачиваем ${msg.total} треков...`;
      logAdd('info', `▶ Начало: ${msg.total} треков`);
      if (msg.track_list?.length) {
        renderTracklistPreview(msg.track_list);
      }
      break;

    case 'track':
      setProgress(msg.done, msg.total);
      dom.progressLabel.textContent = msg.query;
      if (msg.ok) {
        const cls = msg.source === 'ytdlp' ? 'ok-yt' : 'ok';
        const tag = msg.source === 'ytdlp' ? ' [yt]' : '';
        logAdd(cls, `✓ [${msg.done}/${msg.total}]${tag} ${msg.query}`);
      } else {
        logAdd('fail', `✗ [${msg.done}/${msg.total}] ${msg.query}  → ${msg.error}`);
      }
      break;

    case 'done':
      setProgress(msg.total || 1, msg.total || 1);
      dom.progressLabel.textContent = '✅ Готово!';
      dom.statOk.textContent   = msg.ok;
      dom.statFail.textContent = msg.fail;
      dom.statPct.textContent  = msg.pct + '%';
      dom.stats.hidden = false;
      dom.downloadBtn.href = API.result + msg.job_id;
      dom.downloadBtn.hidden = false;
      dom.resetBtn.hidden = false;
      logAdd('info', `✅ ${msg.ok} скачано, ${msg.fail} не найдено (${msg.pct}%)`);
      break;

    case 'error':
      logAdd('fail', '❌ ' + msg.message);
      break;
  }
}

// ── SSE-ридер ─────────────────────────────────────────────────────────────

async function readSSE(response, onMsg) {
  const reader = response.body.getReader();
  const dec    = new TextDecoder();
  let   buf    = '';

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    buf += dec.decode(value, { stream: true });
    const events = buf.split('\n\n');
    buf = events.pop(); // последний может быть неполным

    for (const event of events) {
      const dataLine = event.split('\n').find(l => l.startsWith('data:'));
      if (!dataLine) continue;
      try {
        onMsg(JSON.parse(dataLine.slice(5)));
      } catch { /* пропускаем битые события */ }
    }
  }
}

// ── Сброс UI ──────────────────────────────────────────────────────────────

function resetUI() {
  dom.progressSection.hidden = true;
  dom.stats.hidden           = true;
  dom.downloadBtn.hidden     = true;
  dom.resetBtn.hidden        = true;
  dom.tracklistWrap.hidden   = true;
  dom.log.innerHTML          = '';
  dom.progressBar.style.width = '0%';
  dom.progressPct.textContent = '0%';
  dom.fileName.textContent    = '';
  dom.fileInput.value         = '';
  dom.submitBtn.disabled      = false;
}

dom.resetBtn.addEventListener('click', resetUI);

// ── Submit ────────────────────────────────────────────────────────────────

dom.form.addEventListener('submit', async e => {
  e.preventDefault();

  if (activeTab === 'txt' && !dom.fileInput.files[0]) {
    alert('Выбери .txt файл');
    return;
  }
  if (activeTab === 'yandex' && !dom.yandexUrl.value.trim()) {
    alert('Вставь ссылку на плейлист');
    return;
  }

  let url, init;

  if (activeTab === 'txt') {
    const fd = new FormData();
    fd.append('file',    dom.fileInput.files[0]);
    fd.append('workers', dom.optWorkers.value);
    fd.append('limit',   dom.optLimit.value);
    url  = API.streamTxt;
    init = { method: 'POST', body: fd };
  } else {
    url  = API.streamYandex;
    init = {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify({
        url:   dom.yandexUrl.value.trim(),
        token: dom.yandexToken.value.trim(),
      }),
    };
  }

  dom.submitBtn.disabled      = true;
  dom.log.innerHTML           = '';
  dom.stats.hidden            = true;
  dom.downloadBtn.hidden      = true;
  dom.resetBtn.hidden         = true;
  dom.progressBar.style.width = '0%';
  dom.progressPct.textContent = '0%';
  dom.progressLabel.textContent = 'Подключение...';
  dom.progressSection.hidden  = false;

  try {
    const resp = await fetch(url, init);

    if (!resp.ok) {
      const text = await resp.text();
      logAdd('fail', `❌ Ошибка ${resp.status}: ${text}`);
      dom.submitBtn.disabled = false;
      return;
    }

    await readSSE(resp, handleMsg);
  } catch (err) {
    logAdd('fail', '❌ Сетевая ошибка: ' + err.message);
  }

  dom.submitBtn.disabled = false;
});

// ── Утилиты ───────────────────────────────────────────────────────────────

function escHtml(s) {
  return s.replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;');
}

window.switchTab = switchTab;
