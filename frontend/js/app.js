/* CV Generator — frontend logic */
'use strict';

const API_URL = '/api/generate';

// ── DOM references ────────────────────────────────────────
const form           = document.getElementById('cv-form');
const templateSel    = document.getElementById('template');
const japanFields    = document.getElementById('japan-fields');
const japanLabel     = document.querySelector('.japan-only.subtitle-badge'); // optional
const japanOnlyEls   = document.querySelectorAll('.japan-only');

const expList        = document.getElementById('experience-list');
const eduList        = document.getElementById('education-list');
const langList       = document.getElementById('language-list');

const tplExp         = document.getElementById('tpl-experience');
const tplEdu         = document.getElementById('tpl-education');
const tplLang        = document.getElementById('tpl-language');

const generateBtn    = document.getElementById('generate-btn');
const btnLabel       = document.getElementById('btn-label');
const btnSpinner     = document.getElementById('btn-spinner');
const errorMsg       = document.getElementById('error-msg');

// ── Template toggle ───────────────────────────────────────
templateSel.addEventListener('change', () => {
  const isJapan = templateSel.value === 'japan';
  japanOnlyEls.forEach(el => el.classList.toggle('hidden', !isJapan));
});

// ── Add / remove dynamic items ────────────────────────────
function addItem(list, tpl) {
  const clone = tpl.content.cloneNode(true);
  const block = clone.querySelector('.item-block');
  block.querySelector('.btn-remove').addEventListener('click', () => block.remove());
  list.appendChild(clone);
}

document.getElementById('add-experience').addEventListener('click',
  () => addItem(expList, tplExp));
document.getElementById('add-education').addEventListener('click',
  () => addItem(eduList, tplEdu));
document.getElementById('add-language').addEventListener('click',
  () => addItem(langList, tplLang));

// ── Collect form data ─────────────────────────────────────
function collectData() {
  const get = id => document.getElementById(id).value.trim();

  const experience = Array.from(expList.querySelectorAll('.item-block')).map(b => ({
    position:    b.querySelector('[name="position"]').value.trim(),
    company:     b.querySelector('[name="company"]').value.trim(),
    start_date:  b.querySelector('[name="start_date"]').value.trim(),
    end_date:    b.querySelector('[name="end_date"]').value.trim(),
    description: b.querySelector('[name="description"]').value.trim(),
  })).filter(e => e.company || e.position);

  const education = Array.from(eduList.querySelectorAll('.item-block')).map(b => ({
    institution: b.querySelector('[name="institution"]').value.trim(),
    degree:      b.querySelector('[name="degree"]').value.trim(),
    field:       b.querySelector('[name="field"]').value.trim(),
    start_date:  b.querySelector('[name="start_date"]').value.trim(),
    end_date:    b.querySelector('[name="end_date"]').value.trim(),
  })).filter(e => e.institution);

  const languages = Array.from(langList.querySelectorAll('.item-block')).map(b => ({
    name:        b.querySelector('[name="lang_name"]').value.trim(),
    proficiency: b.querySelector('[name="proficiency"]').value.trim(),
  })).filter(l => l.name);

  const skillsRaw = get('skills');
  const skills = skillsRaw ? skillsRaw.split(',').map(s => s.trim()).filter(Boolean) : [];

  return {
    name:        get('name'),
    email:       get('email'),
    phone:       get('phone'),
    address:     get('address'),
    website:     get('website'),
    birth_date:  get('birth_date'),
    gender:      document.getElementById('gender').value,
    nationality: get('nationality'),
    summary:     get('summary'),
    experience,
    education,
    skills,
    languages,
    template:    templateSel.value,
    format:      document.getElementById('format').value,
  };
}

// ── Submit ────────────────────────────────────────────────
form.addEventListener('submit', async (e) => {
  e.preventDefault();
  hideError();

  const data = collectData();

  if (!data.name) {
    showError('Please enter your full name.');
    document.getElementById('name').focus();
    return;
  }

  setLoading(true);

  try {
    const resp = await fetch(API_URL, {
      method:  'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify(data),
    });

    if (!resp.ok) {
      const msg = await resp.text();
      throw new Error(msg || `Server error (${resp.status})`);
    }

    // Trigger file download
    const blob = await resp.blob();
    const ext  = data.format === 'word' ? 'docx' : 'pdf';
    const url  = URL.createObjectURL(blob);
    const a    = document.createElement('a');
    a.href     = url;
    a.download = `cv_${sanitizeFilename(data.name)}.${ext}`;
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);

  } catch (err) {
    showError('Failed to generate CV: ' + err.message);
  } finally {
    setLoading(false);
  }
});

// ── Helpers ───────────────────────────────────────────────
function setLoading(on) {
  generateBtn.disabled = on;
  btnLabel.textContent = on ? 'Generating…' : 'Generate CV';
  btnSpinner.classList.toggle('hidden', !on);
}

function showError(msg) {
  errorMsg.textContent = msg;
  errorMsg.classList.remove('hidden');
  errorMsg.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}

function hideError() {
  errorMsg.classList.add('hidden');
}

function sanitizeFilename(name) {
  return name.replace(/[^a-zA-Z0-9\u3000-\u9FFF_-]/g, '_').slice(0, 40);
}
