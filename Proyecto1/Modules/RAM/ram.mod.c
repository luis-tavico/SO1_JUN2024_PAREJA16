#include <linux/module.h>
#define INCLUDE_VERMAGIC
#include <linux/build-salt.h>
#include <linux/elfnote-lto.h>
#include <linux/vermagic.h>
#include <linux/compiler.h>

BUILD_SALT;
BUILD_LTO_INFO;

MODULE_INFO(vermagic, VERMAGIC_STRING);
MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};

#ifdef CONFIG_RETPOLINE
MODULE_INFO(retpoline, "Y");
#endif

static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0x602c1205, "module_layout" },
	{ 0xacb46524, "seq_read" },
	{ 0x363f205f, "remove_proc_entry" },
	{ 0x92997ed8, "_printk" },
	{ 0x4334b731, "proc_create" },
	{ 0x8d286714, "seq_printf" },
	{ 0x40c7247c, "si_meminfo" },
	{ 0x5b8239ca, "__x86_return_thunk" },
	{ 0x2dcf0dc0, "single_open" },
	{ 0xbdfb6dbb, "__fentry__" },
};

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "65366731F4DBFE1302BFD74");
