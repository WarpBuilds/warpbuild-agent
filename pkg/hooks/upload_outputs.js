// Uploads a Claude sandbox session's deliverables to backend-cache using the WarpBuild cache client,
// keyed by the anthropic_session_id. Driven by the agentd outputs-upload hook.
//
// Env: WARPBUILD_CACHE_URL, WARPBUILD_RUNNER_VERIFICATION_TOKEN (warp-cache auth), SANDBOX_OUTPUTS_DIR,
// SANDBOX_SESSION_ID. GITHUB_* are unset for a sandbox; the cache service resolves the runner as a
// claude_agent and skips the CI/VCS requirements.
const cache = require("@warpbuilds/cache");

async function main() {
	const outputsDir = process.env.SANDBOX_OUTPUTS_DIR;
	const key = process.env.SANDBOX_SESSION_ID;
	if (!outputsDir || !key) {
		console.error("upload-outputs: SANDBOX_OUTPUTS_DIR and SANDBOX_SESSION_ID are required");
		process.exit(2);
	}
	try {
		const cacheId = await cache.saveCache([outputsDir], key);
		console.log(`upload-outputs: uploaded session ${key} (cache id ${cacheId})`);
	} catch (err) {
		console.error(`upload-outputs: saveCache failed: ${err && err.message ? err.message : err}`);
		process.exit(1);
	}
}

main();
