<script lang="ts">
	const currentYear = new Date().getFullYear();
</script>

<footer class="border-t border-border bg-card mt-auto">
	<div class="container mx-auto px-4 py-6">
		<div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
			<div class="flex items-center gap-3">
				<img src="/logo-red.png" alt="ilionx" class="h-6 w-auto opacity-70" />
				<p class="text-sm text-muted-foreground">
					&copy; {currentYear} ilionx. <span class="font-light italic">experts in eenvoud</span>
				</p>
			</div>
			<p class="text-xs text-muted-foreground">
				Demo application built with SvelteKit & Spring Boot
			</p>
		</div>
	</div>
</footer>
